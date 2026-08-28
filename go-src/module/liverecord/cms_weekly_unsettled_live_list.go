package liverecord

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtime"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liverecorddto"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

// GetCMSWeeklyUnsettledLiveList CMS分页查询本周未结算流水
func GetCMSWeeklyUnsettledLiveList(_ context.Context, req *liverecorddto.CMSWeeklyUnsettledLiveListReq) (*httpserver.CMSQueryResp, error) {
	weekStart, weekEnd := xrtime.WeekDateRange(time.Now())
	roomIds := liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds)
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
		RoomIds:       roomIds,
		LiveDateStart: weekStart,
		LiveDateEnd:   weekEnd,
		Keyword:       req.Keyword,
		Settled:       0,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	profileMap := userinfodao.GetUserProfileMapByUserIds(collectWeeklyUnsettledRoomIds(rows))
	weeklyIncomeMap := liveroomdao.ListWeeklyUnsettledIncomeByRoomIds(collectWeeklyUnsettledRoomIds(rows), weekStart, weekEnd)
	list := make([]*liverecorddto.CMSWeeklyUnsettledLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := toCMSWeeklyUnsettledLiveItem(row, profileMap, weeklyIncomeMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func collectWeeklyUnsettledRoomIds(rows []*liveentity.DailyAnchorEffectiveLive) []uint64 {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.RoomId == 0 {
			continue
		}
		if _, ok := seen[row.RoomId]; ok {
			continue
		}
		seen[row.RoomId] = struct{}{}
		ids = append(ids, row.RoomId)
	}
	return ids
}

func toCMSWeeklyUnsettledLiveItem(row *liveentity.DailyAnchorEffectiveLive, profileMap map[uint64]*userentity.UserInfo, weeklyIncomeMap map[uint64]float64) *liverecorddto.CMSWeeklyUnsettledLiveItem {
	if row == nil || row.ID == "" {
		return nil
	}
	item := &liverecorddto.CMSWeeklyUnsettledLiveItem{
		ID:           row.ID,
		RoomId:       row.RoomId,
		LiveDate:     row.LiveDate,
		LiveDuration: row.LiveDuration,
		TotalIncome:  row.TotalIncome,
	}
	if profileMap != nil {
		if profile := profileMap[row.RoomId]; profile != nil {
			item.RoomNickname = profile.Nickname
			item.RoomAvatar = upload.ResolveAvatarUrlForUser(row.RoomId, profile.Avatar)
		}
	}
	if weeklyIncomeMap != nil {
		item.WeeklyUnsettledTotalIncome = weeklyIncomeMap[row.RoomId]
	}
	return item
}
