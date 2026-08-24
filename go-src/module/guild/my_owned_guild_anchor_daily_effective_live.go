package guild

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/guilddto"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

func collectOwnedGuildIdsForCMS(ctx context.Context) ([]uint64, error) {
	cmsUserId, err := getCMSUserId(ctx)
	if err != nil {
		return nil, err
	}
	return CollectOwnedGuildIdsForCMSUser(cmsUserId), nil
}

// CollectOwnedGuildIdsForCMSUser 获取 CMS 用户作为会长管理的工会 ID 列表
func CollectOwnedGuildIdsForCMSUser(cmsUserId uint64) []uint64 {
	if cmsUserId == 0 {
		return nil
	}
	return collectOwnedGuildIds(guilddao.ListGuildsByLeaderId(cmsUserId))
}

func validateMyOwnedGuildAnchorAccess(ctx context.Context, anchorId uint64) error {
	if anchorId == 0 {
		return nil
	}
	guildIds, err := collectOwnedGuildIdsForCMS(ctx)
	if err != nil {
		return err
	}
	room := liveroomdao.ResolveRoom(anchorId)
	if room == nil || !containsGuildId(guildIds, room.GuildId) {
		return errercode.CreateCode(errercode.NoPermission)
	}
	return nil
}

// GetMyOwnedGuildAnchorList CMS分页查询当前用户作为会长管理的全部工会名下主播
func GetMyOwnedGuildAnchorList(ctx context.Context, req *guilddto.GetMyOwnedGuildAnchorListReq) (*httpserver.CMSQueryResp, error) {
	guildIds, err := collectOwnedGuildIdsForCMS(ctx)
	if err != nil {
		return nil, err
	}
	if len(guildIds) == 0 {
		return httpserver.NewCMSQueryResp(0, []*guilddto.MyGuildAnchorListItem{}), nil
	}
	resp, err := liveroom.QueryAnchorListByGuildIds(ctx, guildIds, &accountdto.QueryAnchorListReq{
		CMSQueryReq: req.CMSQueryReq,
		Key:         req.Key,
	})
	if err != nil {
		return nil, err
	}
	items, _ := resp.Data.([]*accountdto.AnchorListItem)
	list := make([]*guilddto.MyGuildAnchorListItem, 0, len(items))
	for _, item := range items {
		if row := toMyGuildAnchorListItem(item); row != nil {
			list = append(list, row)
		}
	}
	return httpserver.NewCMSQueryResp(resp.Total, list), nil
}

// GetMyOwnedGuildAnchorDailyEffectiveLiveList CMS分页查询当前用户作为会长管理的工会名下主播每日流水
func GetMyOwnedGuildAnchorDailyEffectiveLiveList(ctx context.Context, req *guilddto.CMSMyGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guildIds, err := collectOwnedGuildIdsForCMS(ctx)
	if err != nil {
		return nil, err
	}
	if len(guildIds) == 0 {
		return httpserver.NewCMSQueryResp(0, []*guilddto.GuildAnchorDailyEffectiveLiveItem{}), nil
	}
	roomId := parseGuildIdFilter(req.RoomId)
	if err = validateMyOwnedGuildAnchorAccess(ctx, roomId); err != nil {
		return nil, err
	}
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIds(&liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIdsFilter{
		GuildIds:      guildIds,
		RoomId:        roomId,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	})
	roomIds := collectDailyAnchorRoomIds(rows)
	nicknameMap := userinfodao.GetNicknameMapByUserIds(roomIds)
	list := make([]*guilddto.GuildAnchorDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := toGuildAnchorDailyEffectiveLiveItem(row, nicknameMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
