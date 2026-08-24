package guild

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/guilddto"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// GetMyGuildAnchorDailyEffectiveLiveList CMS分页查询当前用户作为会长管理的工会名下指定主播每日流水
func GetMyGuildAnchorDailyEffectiveLiveList(ctx context.Context, req *guilddto.GetMyGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	if err := validateMyGuildAnchorAccess(ctx, req.GuildId, req.AnchorId); err != nil {
		return nil, err
	}
	return liveroom.QueryAnchorDailyEffectiveLiveList(ctx, &accountdto.GetAnchorDailyEffectiveLiveListReq{
		CMSQueryReq:   req.CMSQueryReq,
		AnchorId:      req.AnchorId,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
		Settled:       req.Settled,
	})
}

func validateMyGuildAnchorAccess(ctx context.Context, guildId, anchorId uint64) error {
	if guildId == 0 || anchorId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := getGuildOwnedByCMSUser(ctx, guildId); err != nil {
		return err
	}
	room := liveroomdao.ResolveRoom(anchorId)
	if room == nil || room.GuildId != guildId {
		return errercode.CreateCode(errercode.NoPermission)
	}
	return nil
}
