package game

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/gamebetdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/upload"
)

const (
	gameConsumeRankTickInterval = 3 * time.Minute
	gameConsumeRankRefreshDelay = 10 * time.Minute
)

type gameConsumeRankRow struct {
	UserId      uint64
	TotalAmount float64
}

type gameConsumeRankSnapshot struct {
	Rows      []*gameConsumeRankRow
	UpdatedAt int64
}

var (
	gameConsumeRankCache           = gmap.NewKVMap[uint64, *gameConsumeRankSnapshot](true)
	gameConsumeRankRefreshDeadline = gmap.NewKVMap[uint64, int64](true)
)

func initGameConsumeRank() {
	for _, roomId := range liveroomdao.ListLivingRoomIds() {
		room := liveroomdao.GetRoomById(roomId)
		if room != nil && room.LiveRecordId > 0 {
			loadGameConsumeRankCache(room.LiveRecordId)
		}
	}
	event.Sub(gameevent.GameBetCreatedEvent, onGameConsumeRankBetCreatedEvent)
	event.Sub(gameevent.RankListRefreshEvent, onGameConsumeRankRefreshEvent)
	xrtimer.AddSingleton(gctx.New(), gameConsumeRankTickInterval, func(ctx context.Context) {
		tryRefreshGameConsumeRankCaches()
	})
}

func onGameConsumeRankBetCreatedEvent(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.LiveRecordId == 0 || ev.Amount <= 0 {
		return
	}
	markGameConsumeRankDataChanged(ev.LiveRecordId)
}

func onGameConsumeRankRefreshEvent(_ any) {
	for _, liveRecordId := range gameConsumeRankCache.Keys() {
		markGameConsumeRankDataChanged(liveRecordId)
	}
	for _, roomId := range liveroomdao.ListLivingRoomIds() {
		room := liveroomdao.GetRoomById(roomId)
		if room != nil && room.LiveRecordId > 0 {
			markGameConsumeRankDataChanged(room.LiveRecordId)
		}
	}
}

func markGameConsumeRankDataChanged(liveRecordId uint64) {
	if liveRecordId == 0 {
		return
	}
	gameConsumeRankRefreshDeadline.Set(liveRecordId, time.Now().Add(gameConsumeRankRefreshDelay).Unix())
}

func tryRefreshGameConsumeRankCaches() {
	nowUnix := time.Now().Unix()
	for _, liveRecordId := range gameConsumeRankRefreshDeadline.Keys() {
		deadlineUnix := gameConsumeRankRefreshDeadline.Get(liveRecordId)
		if deadlineUnix == 0 {
			gameConsumeRankRefreshDeadline.Remove(liveRecordId)
			continue
		}
		if nowUnix > deadlineUnix {
			gameConsumeRankRefreshDeadline.Remove(liveRecordId)
			continue
		}
		loadGameConsumeRankCache(liveRecordId)
	}
}

func loadGameConsumeRankCache(liveRecordId uint64) {
	if liveRecordId == 0 {
		return
	}
	rows := loadGameConsumeRankRows(liveRecordId)
	gameConsumeRankCache.Set(liveRecordId, &gameConsumeRankSnapshot{
		Rows:      rows,
		UpdatedAt: time.Now().Unix(),
	})
}

func loadGameConsumeRankRows(liveRecordId uint64) []*gameConsumeRankRow {
	statRows := gamebetdao.SumGameBetByLiveRecord(liveRecordId)
	list := make([]*gameConsumeRankRow, 0, len(statRows))
	for _, row := range statRows {
		if row == nil || row.UserId == 0 {
			continue
		}
		list = append(list, &gameConsumeRankRow{
			UserId:      row.UserId,
			TotalAmount: row.TotalAmount,
		})
	}
	return list
}

func getGameConsumeRankSnapshot(liveRecordId uint64) *gameConsumeRankSnapshot {
	return gameConsumeRankCache.Get(liveRecordId)
}

func normalizeGameConsumeRankPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func calcUserAge(birthday *time.Time) int {
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

func findMyGameConsumeRank(rows []*gameConsumeRankRow, userId uint64) (rank int, amount float64) {
	rank = -1
	if userId == 0 {
		return rank, 0
	}
	for i, row := range rows {
		if row != nil && row.UserId == userId {
			return i + 1, row.TotalAmount
		}
	}
	return rank, 0
}

// GetAppGameConsumeRank App 分页查询单场直播游戏消费榜
func GetAppGameConsumeRank(ctx context.Context, req *gamebetdto.AppGameConsumeRankReq) (*gamebetdto.AppGameConsumeRankRes, error) {
	liveRecord := liveroomdao.GetLiveRecordById(req.LiveRecordId)
	if liveRecord == nil || liveRecord.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	requestUserId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizeGameConsumeRankPage(req.Page, req.PageSize)
	snapshot := getGameConsumeRankSnapshot(req.LiveRecordId)
	if snapshot == nil {
		loadGameConsumeRankCache(req.LiveRecordId)
		snapshot = getGameConsumeRankSnapshot(req.LiveRecordId)
	}
	rows := snapshot.Rows
	if rows == nil {
		rows = make([]*gameConsumeRankRow, 0)
	}

	total := len(rows)
	myRank, myConsumeAmount := findMyGameConsumeRank(rows, requestUserId)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	list := make([]*gamebetdto.AppGameConsumeRankItem, 0, end-start)
	for i := start; i < end; i++ {
		row := rows[i]
		if row == nil || row.UserId == 0 {
			continue
		}
		item := &gamebetdto.AppGameConsumeRankItem{
			Rank:          i + 1,
			UserId:        strconv.FormatUint(row.UserId, 10),
			ConsumeAmount: row.TotalAmount,
		}
		if u := userinfodao.GetUserInfoByUserId(row.UserId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrlForUser(row.UserId, u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcUserAge(u.Birthday)
		}
		list = append(list, item)
	}

	updatedAt := int64(0)
	if snapshot != nil {
		updatedAt = snapshot.UpdatedAt
	}

	return &gamebetdto.AppGameConsumeRankRes{
		LiveRecordId:  strconv.FormatUint(req.LiveRecordId, 10),
		MyRank:        myRank,
		ConsumeAmount: myConsumeAmount,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		UpdatedAt:     updatedAt,
		List:          list,
	}, nil
}
