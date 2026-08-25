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
	return validateMyOwnedGuildAnchorsAccess(ctx, []uint64{anchorId})
}

func validateMyOwnedGuildAnchorsAccess(ctx context.Context, anchorIds []uint64) error {
	if len(anchorIds) == 0 {
		return nil
	}
	guildIds, err := collectOwnedGuildIdsForCMS(ctx)
	if err != nil {
		return err
	}
	for _, anchorId := range anchorIds {
		if anchorId == 0 {
			continue
		}
		room := liveroomdao.ResolveRoom(anchorId)
		if room == nil || !containsGuildId(guildIds, room.GuildId) {
			return errercode.CreateCode(errercode.NoPermission)
		}
	}
	return nil
}

func resolveOwnedGuildAnchorRoomIds(roomId string, roomIds []string) []uint64 {
	if len(roomIds) > 0 {
		return parseGuildIdFilters(roomIds)
	}
	if id := parseGuildIdFilter(roomId); id > 0 {
		return []uint64{id}
	}
	return nil
}

// ResolveOwnedGuildAnchorRoomIds 解析主播筛选 ID 列表
func ResolveOwnedGuildAnchorRoomIds(roomId string, roomIds []string) []uint64 {
	return resolveOwnedGuildAnchorRoomIds(roomId, roomIds)
}

// ValidateMyOwnedGuildAnchorsAccessForUser 校验 CMS 用户是否可访问这些主播
func ValidateMyOwnedGuildAnchorsAccessForUser(cmsUserId uint64, anchorIds []uint64) error {
	if len(anchorIds) == 0 {
		return nil
	}
	guildIds := CollectOwnedGuildIdsForCMSUser(cmsUserId)
	for _, anchorId := range anchorIds {
		if anchorId == 0 {
			continue
		}
		room := liveroomdao.ResolveRoom(anchorId)
		if room == nil || !containsGuildId(guildIds, room.GuildId) {
			return errercode.CreateCode(errercode.NoPermission)
		}
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
	roomIds := resolveOwnedGuildAnchorRoomIds(req.RoomId, req.RoomIds)
	if err = validateMyOwnedGuildAnchorsAccess(ctx, roomIds); err != nil {
		return nil, err
	}
	filter := &liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIdsFilter{
		GuildIds:      guildIds,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
		PageIndex:     req.PageIndex,
		PageSize:      req.PageSize,
	}
	if len(roomIds) == 1 {
		filter.RoomId = roomIds[0]
	} else if len(roomIds) > 1 {
		filter.RoomIds = roomIds
	}
	total, rows := liveroomdao.DailyAnchorEffectiveLiveCMSListByGuildIds(filter)
	anchorRoomIds := CollectDailyAnchorRoomIds(rows)
	profileMap := userinfodao.GetUserProfileMapByUserIds(anchorRoomIds)
	list := make([]*guilddto.GuildAnchorDailyEffectiveLiveItem, 0, len(rows))
	for _, row := range rows {
		if item := ToGuildAnchorDailyEffectiveLiveItem(row, profileMap); item != nil {
			list = append(list, item)
		}
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}
