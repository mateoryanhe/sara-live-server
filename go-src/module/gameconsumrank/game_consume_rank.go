package gameconsumrank

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gameconsumrankdto"
	"xr-game-server/gameevent"
	"xr-game-server/module/upload"
)

const (
	defaultPageSize             = 20
	maxPageSize                 = 100
	gameConsumeRankTickInterval = 3 * time.Minute
	gameConsumeRankRefreshDelay = 10 * time.Minute
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

var gameConsumeRankCache atomic.Value
var dataRefreshDeadline atomic.Int64

func init() {
	gameConsumeRankCache.Store(&rankSnapshot{
		Today:  make([]*rankItem, 0),
		Last7:  make([]*rankItem, 0),
		Last30: make([]*rankItem, 0),
	})
}

// Init 初始化游戏消费榜缓存,订阅游戏下注事件,每3分钟检查是否需要刷新
func Init() {
	loadGameConsumeRankCache()
	event.Sub(gameevent.GameBetCreatedEvent, onGameBetCreatedEvent)
	event.Sub(gameevent.RankListRefreshEvent, onRankListRefreshEvent)
	xrtimer.AddSingleton(gctx.New(), gameConsumeRankTickInterval, func(ctx context.Context) {
		tryRefreshGameConsumeRankCache()
	})
}

func onGameBetCreatedEvent(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.Amount <= 0 {
		return
	}
	markGameConsumeRankDataChanged()
}

func markGameConsumeRankDataChanged() {
	dataRefreshDeadline.Store(time.Now().Add(gameConsumeRankRefreshDelay).Unix())
}

func onRankListRefreshEvent(_ any) {
	markGameConsumeRankDataChanged()
}

func tryRefreshGameConsumeRankCache() {
	deadlineUnix := dataRefreshDeadline.Load()
	if deadlineUnix == 0 {
		return
	}
	now := time.Now()
	if now.After(time.Unix(deadlineUnix, 0)) {
		dataRefreshDeadline.Store(0)
		return
	}
	loadGameConsumeRankCache()
}

func loadGameConsumeRankCache() {
	now := time.Now()
	snapshot := &rankSnapshot{
		Today:     buildRankItems(gamebetdao.SumGameBetByUser(startOfDay(now), now)),
		Last7:     buildRankItems(gamebetdao.SumGameBetByUser(now.AddDate(0, 0, -7), now)),
		Last30:    buildRankItems(gamebetdao.SumGameBetByUser(now.AddDate(0, 0, -30), now)),
		UpdatedAt: now.Unix(),
	}
	gameConsumeRankCache.Store(snapshot)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func buildRankItems(rows []*gamebetdao.GameConsumeRankRow) []*rankItem {
	if len(rows) == 0 {
		return make([]*rankItem, 0)
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
			ConsumeAmount: row.TotalAmount,
		}
		if profile := userinfodao.GetUserInfoByUserId(row.UserId); profile != nil {
			item.Nickname = profile.Nickname
			item.Avatar = upload.ResolveAvatarUrlForUser(row.UserId, profile.Avatar)
			item.VipLevel = profile.VipLevel
			item.Gender = profile.Gender
			item.Age = calcAge(profile.Birthday)
		}
		list = append(list, item)
	}
	return list
}

func getSnapshot() *rankSnapshot {
	v := gameConsumeRankCache.Load()
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
	case gameconsumrankdto.GameConsumeRankPeriodToday:
		return snapshot.Today
	case gameconsumrankdto.GameConsumeRankPeriodLast7:
		return snapshot.Last7
	case gameconsumrankdto.GameConsumeRankPeriodLast30:
		return snapshot.Last30
	default:
		return make([]*rankItem, 0)
	}
}

// GetAppGameConsumeRankList App端分页查询游戏消费榜
func GetAppGameConsumeRankList(ctx context.Context, req *gameconsumrankdto.AppGameConsumeRankListReq) (*gameconsumrankdto.AppGameConsumeRankListRes, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	snapshot := getSnapshot()
	all := getRankListByPeriod(snapshot, req.Period)
	total := len(all)
	myRank := findMyRank(all, httpserver.GetAuthId(ctx))
	start, end := pageRange(total, page, pageSize)
	pageData := make([]*gameconsumrankdto.AppGameConsumeRankItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		pageData = append(pageData, &gameconsumrankdto.AppGameConsumeRankItem{
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
	return &gameconsumrankdto.AppGameConsumeRankListRes{
		Period:    req.Period,
		MyRank:    myRank,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		UpdatedAt: snapshot.UpdatedAt,
		List:      pageData,
	}, nil
}

func findMyRank(rows []*rankItem, userId uint64) int {
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
