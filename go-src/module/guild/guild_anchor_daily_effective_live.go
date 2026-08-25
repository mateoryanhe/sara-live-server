package guild

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

// GetGuildAnchorDailyEffectiveLiveList CMS分页查询指定工会名下主播每日流水
func GetGuildAnchorDailyEffectiveLiveList(_ context.Context, req *guilddto.CMSGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guildId := parseGuildIdFilter(req.GuildId)
	if guildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSListByGuild(&liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildFilter{
		GuildId:       guildId,
		RoomId:        parseGuildIdFilter(req.RoomId),
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	roomIds := CollectDailyAnchorRoomIds(rows)
	profileMap := userinfodao.GetUserProfileMapByUserIds(roomIds)
	unsettledIncomeMap := liveroomdao.ListLiveRoomIncomeUnsettledTotalForCMS(roomIds)
	list := make([]*guilddto.GuildAnchorDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := ToGuildAnchorDailyEffectiveLiveItem(row, profileMap, unsettledIncomeMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

func CollectDailyAnchorRoomIds(rows []*liveentity.DailyAnchorEffectiveLive) []uint64 {
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

func ToGuildAnchorDailyEffectiveLiveItem(row *liveentity.DailyAnchorEffectiveLive, profileMap map[uint64]*userentity.UserInfo, unsettledIncomeMap map[uint64]float64) *guilddto.GuildAnchorDailyEffectiveLiveItem {
	if row == nil || row.ID == "" {
		return nil
	}
	item := &guilddto.GuildAnchorDailyEffectiveLiveItem{
		ID:                        row.ID,
		RoomId:                    row.RoomId,
		LiveDate:                  row.LiveDate,
		LiveDuration:              row.LiveDuration,
		Settled:                   row.Settled,
		LiveRoomIncomeAmountsItem: toIncomeAmountsItem(&row.LiveRoomIncomeAmounts),
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
