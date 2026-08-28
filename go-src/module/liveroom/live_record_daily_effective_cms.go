package liveroom

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/liverecorddto"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/upload"
)

// GetCMSDailyEffectiveLiveList CMS分页查询每日流水
func GetCMSDailyEffectiveLiveList(_ context.Context, req *liverecorddto.CMSDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	roomIds := liveroomdao.ParseLiveRecordAnchorIds(req.AnchorId, req.PlatformAnchorId, req.GuildAnchorId, req.AnchorIds)
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSMultiList(&liveroomdao.DailyAnchorEffectiveLiveCMSMultiListFilter{
		RoomIds:       roomIds,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	profileMap := userinfodao.GetUserProfileMapByUserIds(collectDailyAnchorRoomIds(rows))
	unsettledIncomeMap := liveroomdao.ListLiveRoomIncomeUnsettledTotalForCMS(collectDailyAnchorRoomIds(rows))
	list := make([]*liverecorddto.CMSDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := toCMSDailyEffectiveLiveItem(row, profileMap, unsettledIncomeMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func collectDailyAnchorRoomIds(rows []*liveentity.DailyAnchorEffectiveLive) []uint64 {
	if len(rows) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.RoomId > 0 {
			ids = append(ids, row.RoomId)
		}
	}
	return ids
}

func toCMSDailyEffectiveLiveItem(row *liveentity.DailyAnchorEffectiveLive, profileMap map[uint64]*userentity.UserInfo, unsettledIncomeMap map[uint64]float64) *liverecorddto.CMSDailyEffectiveLiveItem {
	if row == nil || row.ID == "" {
		return nil
	}
	amounts := row.LiveRoomIncomeAmounts
	item := &liverecorddto.CMSDailyEffectiveLiveItem{
		ID:           row.ID,
		RoomId:       row.RoomId,
		LiveDate:     row.LiveDate,
		LiveDuration: row.LiveDuration,
		Settled:      row.Settled,
		LiveRoomIncomeAmountsItem: accountdto.LiveRoomIncomeAmountsItem{
			TotalIncome:                  amounts.TotalIncome,
			TotalGiftIncome:              amounts.TotalGiftIncome,
			TotalPaidDanmakuIncome:       amounts.TotalPaidDanmakuIncome,
			TotalPrivateRoomTicketIncome: amounts.TotalPrivateRoomTicketIncome,
			TotalPrivateRoomWatchIncome:  amounts.TotalPrivateRoomWatchIncome,
			TotalVideoCallIncome:         amounts.TotalVideoCallIncome,
			TotalVideoCallTicketIncome:   amounts.TotalVideoCallTicketIncome,
			TotalVideoCallBillingIncome:  amounts.TotalVideoCallBillingIncome,
			TotalLiveDuration:            amounts.TotalLiveDuration,
		},
	}
	if profileMap != nil {
		if profile := profileMap[row.RoomId]; profile != nil {
			item.RoomNickname = profile.Nickname
			item.RoomAvatar = upload.ResolveAvatarUrlForUser(row.RoomId, profile.Avatar)
		}
	}
	if unsettledIncomeMap != nil {
		item.UnsettledTotalIncome = unsettledIncomeMap[row.RoomId]
	}
	if !row.CreatedAt.IsZero() {
		createdAt := row.CreatedAt
		item.CreatedAt = &createdAt
	}
	if !row.UpdatedAt.IsZero() {
		updatedAt := row.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}
