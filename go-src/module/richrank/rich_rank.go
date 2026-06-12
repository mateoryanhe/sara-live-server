package richrank

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/currencylogdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/richrankdto"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	defaultPageSize         = 20
	maxPageSize             = 100
	richRankRefreshInterval = 10 * time.Minute
)

type rankItem struct {
	Rank          int
	UserId        uint64
	Nickname      string
	Avatar        string
	ConsumeAmount float64
	VipLevel      uint32
	Gender        uint8
	Age           int
}

type rankSnapshot struct {
	Today     []*rankItem
	Last7     []*rankItem
	Last30    []*rankItem
	UpdatedAt int64
}

var richRankCache atomic.Value

func init() {
	richRankCache.Store(&rankSnapshot{
		Today:  make([]*rankItem, 0),
		Last7:  make([]*rankItem, 0),
		Last30: make([]*rankItem, 0),
	})
}

// Init 初始化富豪榜缓存,并每10分钟刷新一次
func Init() {
	loadRichRankCache()
	xrtimer.AddSingleton(gctx.New(), richRankRefreshInterval, func(ctx context.Context) {
		loadRichRankCache()
	})
}

func loadRichRankCache() {
	now := time.Now()
	snapshot := &rankSnapshot{
		Today:     buildRankItems(currencylogdao.SumDiamondConsumeByUser(startOfDay(now), now)),
		Last7:     buildRankItems(currencylogdao.SumDiamondConsumeByUser(now.AddDate(0, 0, -7), now)),
		Last30:    buildRankItems(currencylogdao.SumDiamondConsumeByUser(now.AddDate(0, 0, -30), now)),
		UpdatedAt: now.Unix(),
	}
	richRankCache.Store(snapshot)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func buildRankItems(rows []*currencylogdao.DiamondConsumeStatRow) []*rankItem {
	if len(rows) == 0 {
		return make([]*rankItem, 0)
	}
	userIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.UserId == 0 {
			continue
		}
		userIds = append(userIds, row.UserId)
	}

	list := make([]*rankItem, 0, len(rows))
	rankNo := 0
	for _, row := range rows {
		if row == nil || row.UserId == 0 {
			continue
		}
		rankNo++
		item := &rankItem{
			Rank:          rankNo,
			UserId:        row.UserId,
			ConsumeAmount: row.Total,
		}
		if profile := userinfodao.GetUserInfoByUserId(row.UserId); profile != nil {
			item.Nickname = profile.Nickname
			item.Avatar = upload.ResolveAvatarUrl(profile.Avatar)
			item.VipLevel = profile.VipLevel
			item.Gender = profile.Gender
			item.Age = calcAge(profile.Birthday)
		}
		list = append(list, item)
	}
	return list
}

func getSnapshot() *rankSnapshot {
	v := richRankCache.Load()
	if v == nil {
		return &rankSnapshot{
			Today:  make([]*rankItem, 0),
			Last7:  make([]*rankItem, 0),
			Last30: make([]*rankItem, 0),
		}
	}
	snapshot, ok := v.(*rankSnapshot)
	if !ok || snapshot == nil {
		return &rankSnapshot{
			Today:  make([]*rankItem, 0),
			Last7:  make([]*rankItem, 0),
			Last30: make([]*rankItem, 0),
		}
	}
	return snapshot
}

func getRankListByPeriod(snapshot *rankSnapshot, period int) []*rankItem {
	switch period {
	case richrankdto.RichRankPeriodToday:
		return snapshot.Today
	case richrankdto.RichRankPeriodLast7:
		return snapshot.Last7
	case richrankdto.RichRankPeriodLast30:
		return snapshot.Last30
	default:
		return make([]*rankItem, 0)
	}
}

// GetAppRichRankList App端分页查询富豪榜
func GetAppRichRankList(_ context.Context, req *richrankdto.AppRichRankListReq) (*richrankdto.AppRichRankListRes, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	snapshot := getSnapshot()
	all := getRankListByPeriod(snapshot, req.Period)
	total := len(all)
	start, end := pageRange(total, page, pageSize)
	pageData := make([]*richrankdto.AppRichRankItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		pageData = append(pageData, &richrankdto.AppRichRankItem{
			Rank:          row.Rank,
			UserId:        strconv.FormatUint(row.UserId, 10),
			Nickname:      row.Nickname,
			Avatar:        row.Avatar,
			ConsumeAmount: row.ConsumeAmount,
			VipLevel:      row.VipLevel,
			Gender:        row.Gender,
			Age:           row.Age,
		})
	}
	return &richrankdto.AppRichRankListRes{
		Period:    req.Period,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		UpdatedAt: snapshot.UpdatedAt,
		List:      pageData,
	}, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

func calcAge(birthday *time.Time) int {
	if birthday == nil || birthday.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	anniversary := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

func pageRange(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}
