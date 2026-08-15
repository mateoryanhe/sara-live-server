package anchorrank

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/anchorrankdto"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
	"xr-game-server/module/upload"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	defaultPageSize        = 20
	maxPageSize            = 100
	anchorRankTickInterval = 3 * time.Minute
	anchorRankRefreshDelay = 10 * time.Minute
)

type rankItem struct {
	Rank          int
	UserId        uint64
	Nickname      string
	Avatar        string
	RevenueAmount uint64
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

var anchorRankCache atomic.Value
var dataRefreshDeadline atomic.Int64

func init() {
	anchorRankCache.Store(&rankSnapshot{
		Today:  make([]*rankItem, 0),
		Last7:  make([]*rankItem, 0),
		Last30: make([]*rankItem, 0),
	})
}

// Init 初始化主播红人榜缓存,订阅收益事件,每3分钟检查是否需要刷新
func Init() {
	loadAnchorRankCache()
	event.Sub(gameevent.RevenueEventEvent, onRevenueEvent)
	event.Sub(gameevent.RankListRefreshEvent, onRankListRefreshEvent)
	xrtimer.AddSingleton(gctx.New(), anchorRankTickInterval, func(ctx context.Context) {
		tryRefreshAnchorRankCache()
	})
}

func onRevenueEvent(data any) {
	log, ok := data.(*entity.LiveRevenueLog)
	if !ok || log == nil {
		return
	}
	if log.ReceiverId == 0 || log.TotalAmount <= 0 {
		return
	}
	if log.Status != entity.LiveRevenueLogStatusNormal {
		return
	}
	markAnchorRankDataChanged()
}

func markAnchorRankDataChanged() {
	dataRefreshDeadline.Store(time.Now().Add(anchorRankRefreshDelay).Unix())
}

func onRankListRefreshEvent(_ any) {
	markAnchorRankDataChanged()
}

func tryRefreshAnchorRankCache() {
	deadlineUnix := dataRefreshDeadline.Load()
	if deadlineUnix == 0 {
		return
	}
	now := time.Now()
	if now.After(time.Unix(deadlineUnix, 0)) {
		dataRefreshDeadline.Store(0)
		return
	}
	loadAnchorRankCache()
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
			item.Avatar = upload.ResolveAvatarUrlForUser(row.ReceiverId, profile.Avatar)
			item.VipLevel = profile.VipLevel
			item.Gender = profile.Gender
			item.Age = calcAge(profile.Birthday)
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
func GetAppAnchorRankList(ctx context.Context, req *anchorrankdto.AppAnchorRankListReq) (*anchorrankdto.AppAnchorRankListRes, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	snapshot := getSnapshot()
	all := getRankListByPeriod(snapshot, req.Period)
	total := len(all)
	myRank := findMyAnchorRank(all, httpserver.GetAuthId(ctx))
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
			VipLevel:      row.VipLevel,
			Gender:        row.Gender,
			Age:           row.Age,
		})
	}
	return &anchorrankdto.AppAnchorRankListRes{
		Period:    req.Period,
		MyRank:    myRank,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		UpdatedAt: snapshot.UpdatedAt,
		List:      pageData,
	}, nil
}

func findMyAnchorRank(rows []*rankItem, userId uint64) int {
	if userId == 0 {
		return -1
	}
	for _, row := range rows {
		if row != nil && row.UserId == userId {
			return row.Rank
		}
	}
	return -1
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

// GetUserLast30DayRevenue 从主播红人榜缓存获取用户最近30天收益,未上榜返回0
func GetUserLast30DayRevenue(userId uint64) uint64 {
	if userId == 0 {
		return 0
	}
	for _, item := range getSnapshot().Last30 {
		if item != nil && item.UserId == userId {
			return item.RevenueAmount
		}
	}
	return 0
}
