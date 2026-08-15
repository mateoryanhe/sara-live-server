package liveroom

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/cmd"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/core/xrpool"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/upload"
)

const (
	contributionRankTickInterval = 3 * time.Minute
	contributionRankRefreshDelay = 10 * time.Minute
)

type contributionRankRow struct {
	SenderId    uint64
	TotalAmount float64
}

type contributionRankSnapshot struct {
	Today     []*contributionRankRow
	Last7     []*contributionRankRow
	Last30    []*contributionRankRow
	UpdatedAt int64
}

var (
	contributionRankCache           = gmap.NewKVMap[uint64, *contributionRankSnapshot](true)
	contributionRankRefreshDeadline = gmap.NewKVMap[uint64, int64](true)
)

func initContributionRank() {
	for _, roomId := range liveroomdao.ListLivingRoomIds() {
		loadContributionRankCache(roomId)
	}
	event.Sub(gameevent.CurrencyChangeEvent, onContributionRankDiamondConsumeEvent)
	event.Sub(gameevent.RankListRefreshEvent, onContributionRankRefreshEvent)
	xrtimer.AddSingleton(gctx.New(), contributionRankTickInterval, func(ctx context.Context) {
		tryRefreshContributionRankCaches()
	})
}

func onContributionRankDiamondConsumeEvent(data any) {
	ev, ok := data.(*gameevent.CurrencyChangeEventData)
	if !ok || ev == nil {
		return
	}
	if ev.Type != gameevent.CurrencyTypeDiamond || ev.Action != gameevent.CurrencyActionSub || ev.Amount <= 0 {
		return
	}
	if ev.Reason != currency.ReasonGiftSend && ev.Reason != currency.ReasonPaidDanmaku {
		return
	}
	for _, roomId := range findOnlineRoomIdsByUser(ev.UserId) {
		markContributionRankDataChanged(roomId)
	}
}

func onContributionRankRefreshEvent(_ any) {
	for _, roomId := range taskMap.Slice() {
		markContributionRankDataChanged(roomId)
	}
}

func markContributionRankDataChanged(roomId uint64) {
	if roomId == 0 {
		return
	}
	contributionRankRefreshDeadline.Set(roomId, time.Now().Add(contributionRankRefreshDelay).Unix())
}

func tryRefreshContributionRankCaches() {
	nowUnix := time.Now().Unix()
	for _, roomId := range contributionRankRefreshDeadline.Keys() {
		// 已下播房间不再刷新,并清掉残留缓存
		if !taskMap.Contains(roomId) {
			clearContributionRankCache(roomId)
			continue
		}
		deadlineUnix := contributionRankRefreshDeadline.Get(roomId)
		if deadlineUnix == 0 {
			contributionRankRefreshDeadline.Remove(roomId)
			continue
		}
		if nowUnix > deadlineUnix {
			contributionRankRefreshDeadline.Remove(roomId)
			continue
		}
		loadContributionRankCache(roomId)
	}
}

func loadContributionRankCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	now := time.Now()
	contributionRankCache.Set(roomId, &contributionRankSnapshot{
		Today:     loadContributionRankRows(roomId, contributionRankStartTime(now, liveroomdto.ContributionRankPeriodToday), now),
		Last7:     loadContributionRankRows(roomId, contributionRankStartTime(now, liveroomdto.ContributionRankPeriodLast7), now),
		Last30:    loadContributionRankRows(roomId, contributionRankStartTime(now, liveroomdto.ContributionRankPeriodLast30), now),
		UpdatedAt: now.Unix(),
	})
}

func clearContributionRankCache(roomId uint64) {
	contributionRankCache.Remove(roomId)
	contributionRankRefreshDeadline.Remove(roomId)
}

func refreshRoomAudienceCaches(roomId uint64) {
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		flushOnlineLists(roomId)
		broadcastAudienceListRefresh(roomId)
	})
}

func broadcastAudienceListRefresh(roomId uint64) {
	online := getOnline(roomId)
	payload := &liveroomdto.AudienceListRefreshPushItem{
		RoomId: roomId,
	}

	for _, userId := range online {
		push.Data(userId, cmd.LiveRoomAudienceListRefresh, payload)
	}
	push.Data(roomId, cmd.LiveRoomAudienceListRefresh, payload)
}

func clearRoomAudienceCaches(roomId uint64) {
	clearRoomOnlineMap(roomId)
	clearOnlineLists(roomId)
	clearContributionRankCache(roomId)
}

func loadContributionRankRows(roomId uint64, startTime, endTime time.Time) []*contributionRankRow {
	list := make([]*contributionRankRow, 0)
	onlineUserIds := getOnline(roomId)
	if len(onlineUserIds) == 0 {
		return list
	}
	rows := liveroomdao.SumAudienceContributionByRoom(roomId, startTime, endTime, onlineUserIds)
	for _, row := range rows {
		if row == nil || row.SenderId == 0 {
			continue
		}
		list = append(list, &contributionRankRow{
			SenderId:    row.SenderId,
			TotalAmount: row.TotalAmount,
		})
	}
	return list
}

func getContributionRankSnapshot(roomId uint64) *contributionRankSnapshot {
	v := contributionRankCache.Get(roomId)
	if v == nil {
		return nil
	}
	return v
}

func getContributionRankRows(snapshot *contributionRankSnapshot, period int) []*contributionRankRow {
	empty := make([]*contributionRankRow, 0)
	if snapshot == nil {
		return empty
	}
	var rows []*contributionRankRow
	switch period {
	case liveroomdto.ContributionRankPeriodToday:
		rows = snapshot.Today
	case liveroomdto.ContributionRankPeriodLast7:
		rows = snapshot.Last7
	case liveroomdto.ContributionRankPeriodLast30:
		rows = snapshot.Last30
	default:
		return empty
	}
	if rows == nil {
		return empty
	}
	return rows
}

// GetContributionRank App端分页查询直播间观众贡献榜(礼物+付费弹幕)
func GetContributionRank(ctx context.Context, req *liveroomdto.GetContributionRankReq) (*liveroomdto.GetContributionRankRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	requestUserId := httpserver.GetAuthId(ctx)
	page, pageSize := normalizeOnlineListPage(req.Page, req.PageSize)
	snapshot := getContributionRankSnapshot(req.RoomId)
	if snapshot == nil {
		loadContributionRankCache(req.RoomId)
		snapshot = getContributionRankSnapshot(req.RoomId)
	}
	rows := getContributionRankRows(snapshot, req.Period)
	if rows == nil {
		rows = make([]*contributionRankRow, 0)
	}

	total := len(rows)
	myRank, myContributionAmount := findMyContribution(rows, requestUserId)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	list := make([]*liveroomdto.ContributionRankItem, 0, end-start)
	for i := start; i < end; i++ {
		row := rows[i]
		if row == nil || row.SenderId == 0 {
			continue
		}
		item := &liveroomdto.ContributionRankItem{
			Rank:               i + 1,
			UserId:             strconv.FormatUint(row.SenderId, 10),
			ContributionAmount: row.TotalAmount,
		}
		if u := userinfodao.GetUserInfoByUserId(row.SenderId); u != nil {
			item.Nickname = u.Nickname
			item.Avatar = upload.ResolveAvatarUrlForUser(row.SenderId, u.Avatar)
			item.VipLevel = u.VipLevel
			item.Gender = u.Gender
			item.Age = calcAge(u.Birthday)
		}
		list = append(list, item)
	}

	updatedAt := int64(0)
	if snapshot != nil {
		updatedAt = snapshot.UpdatedAt
	}

	return &liveroomdto.GetContributionRankRes{
		RoomId:             strconv.FormatUint(req.RoomId, 10),
		Period:             req.Period,
		MyRank:             myRank,
		ContributionAmount: myContributionAmount,
		Total:              total,
		Page:               page,
		PageSize:           pageSize,
		UpdatedAt:          updatedAt,
		List:               list,
	}, nil
}

func findMyContribution(rows []*contributionRankRow, userId uint64) (rank int, amount float64) {
	rank = -1
	if userId == 0 {
		return rank, 0
	}
	for i, row := range rows {
		if row != nil && row.SenderId == userId {
			return i + 1, row.TotalAmount
		}
	}
	return rank, 0
}

func contributionRankStartTime(now time.Time, period int) time.Time {
	switch period {
	case liveroomdto.ContributionRankPeriodToday:
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case liveroomdto.ContributionRankPeriodLast7:
		return now.AddDate(0, 0, -7)
	case liveroomdto.ContributionRankPeriodLast30:
		return now.AddDate(0, 0, -30)
	default:
		return now
	}
}
