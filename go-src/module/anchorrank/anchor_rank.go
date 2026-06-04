package anchorrank

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorrankdto"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	defaultPageSize           = 20
	maxPageSize               = 100
	anchorRankRefreshInterval = 10 * time.Minute
)

type rankItem struct {
	Rank          int
	UserId        uint64
	Nickname      string
	Avatar        string
	RevenueAmount uint64
}

type rankSnapshot struct {
	Today     []*rankItem
	Last7     []*rankItem
	Last30    []*rankItem
	UpdatedAt int64
}

var anchorRankCache atomic.Value

func init() {
	anchorRankCache.Store(&rankSnapshot{
		Today:  make([]*rankItem, 0),
		Last7:  make([]*rankItem, 0),
		Last30: make([]*rankItem, 0),
	})
}

// Init 初始化主播红人榜缓存,并每10分钟刷新一次
func Init() {
	loadAnchorRankCache()
	xrtimer.AddSingleton(gctx.New(), anchorRankRefreshInterval, func(ctx context.Context) {
		loadAnchorRankCache()
	})
}

func loadAnchorRankCache() {
	now := time.Now()
	snapshot := &rankSnapshot{
		Today:     buildRankItems(liveroomdao.SumRevenueByReceiver(startOfDay(now), now)),
		Last7:     buildRankItems(liveroomdao.SumRevenueByReceiver(now.AddDate(0, 0, -7), now)),
		Last30:    buildRankItems(liveroomdao.SumRevenueByReceiver(now.AddDate(0, 0, -30), now)),
		UpdatedAt: now.Unix(),
	}
	anchorRankCache.Store(snapshot)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func buildRankItems(rows []*liveroomdao.AnchorRevenueStatRow) []*rankItem {
	if len(rows) == 0 {
		return make([]*rankItem, 0)
	}
	list := make([]*rankItem, 0, len(rows))
	rankNo := 0
	for _, row := range rows {
		if row == nil || row.ReceiverId == 0 {
			continue
		}
		rankNo++
		item := &rankItem{
			Rank:          rankNo,
			UserId:        row.ReceiverId,
			RevenueAmount: row.TotalAmount,
		}
		if profile := userinfodao.GetUserInfoByUserId(row.ReceiverId); profile != nil {
			item.Nickname = profile.Nickname
			item.Avatar = upload.ResolveAvatarUrl(profile.Avatar)
		}
		list = append(list, item)
	}
	return list
}

func getSnapshot() *rankSnapshot {
	v := anchorRankCache.Load()
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
	case anchorrankdto.AnchorRankPeriodToday:
		return snapshot.Today
	case anchorrankdto.AnchorRankPeriodLast7:
		return snapshot.Last7
	case anchorrankdto.AnchorRankPeriodLast30:
		return snapshot.Last30
	default:
		return make([]*rankItem, 0)
	}
}

// GetAppAnchorRankList App端分页查询主播红人榜
func GetAppAnchorRankList(_ context.Context, req *anchorrankdto.AppAnchorRankListReq) (*anchorrankdto.AppAnchorRankListRes, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	snapshot := getSnapshot()
	all := getRankListByPeriod(snapshot, req.Period)
	total := len(all)
	start, end := pageRange(total, page, pageSize)
	pageData := make([]*anchorrankdto.AppAnchorRankItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		pageData = append(pageData, &anchorrankdto.AppAnchorRankItem{
			Rank:          row.Rank,
			UserId:        strconv.FormatUint(row.UserId, 10),
			Nickname:      row.Nickname,
			Avatar:        row.Avatar,
			RevenueAmount: row.RevenueAmount,
		})
	}
	return &anchorrankdto.AppAnchorRankListRes{
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
