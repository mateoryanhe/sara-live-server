package liveroom

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
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

var contributionRankCache = gmap.NewKVMap[uint64, *contributionRankSnapshot](true)

func flushContributionRankCache(roomId uint64) {
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
}

func refreshRoomAudienceCaches(roomId uint64) {
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		flushContributionRankCache(roomId)
		flushOnlineLists(roomId)
		broadcastAudienceListRefresh(roomId)
	})

}

func broadcastAudienceListRefresh(roomId uint64) {
	online := getOnline(roomId)
	if len(online) == 0 {
		return
	}

	userIds := commonOnlineMap.Get(roomId)
	if userIds == nil {
		userIds = make([]uint64, 0)
	}

	payload := &liveroomdto.AudienceListRefreshPushItem{
		RoomId: roomId,
	}

	for _, userId := range online {
		push.Data(userId, cmd.LiveRoomAudienceListRefresh, payload)
	}
	push.Data(roomId, cmd.LiveRoomAudienceListRefresh, payload)
}

func clearRoomAudienceCaches(roomId uint64) {
	clearOnlineLists(roomId)
	clearContributionRankCache(roomId)
}

func loadContributionRankRows(roomId uint64, startTime, endTime time.Time) []*contributionRankRow {
	onlineUserIds := getOnline(roomId)
	if len(onlineUserIds) == 0 {
		return nil
	}
	rows := liveroomdao.SumAudienceContributionByRoom(roomId, startTime, endTime, onlineUserIds)
	list := make([]*contributionRankRow, 0, len(rows))
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
	if snapshot == nil {
		return nil
	}
	switch period {
	case liveroomdto.ContributionRankPeriodToday:
		return snapshot.Today
	case liveroomdto.ContributionRankPeriodLast7:
		return snapshot.Last7
	case liveroomdto.ContributionRankPeriodLast30:
		return snapshot.Last30
	default:
		return nil
	}
}

// GetContributionRank App端分页查询直播间观众贡献榜(礼物+付费弹幕)
func GetContributionRank(_ context.Context, req *liveroomdto.GetContributionRankReq) (*liveroomdto.GetContributionRankRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	page, pageSize := normalizeOnlineListPage(req.Page, req.PageSize)
	snapshot := getContributionRankSnapshot(req.RoomId)
	rows := getContributionRankRows(snapshot, req.Period)
	if rows == nil {
		rows = make([]*contributionRankRow, 0)
	}

	total := len(rows)
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
			item.Avatar = upload.ResolveAvatarUrl(u.Avatar)
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
		RoomId:    strconv.FormatUint(req.RoomId, 10),
		Period:    req.Period,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		UpdatedAt: updatedAt,
		List:      list,
	}, nil
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
