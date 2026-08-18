package guild

import (
	"context"

	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// SetGuildAnchorType CMS设置工会名下主播类型(仅支持普通/高级主播互切)
func SetGuildAnchorType(ctx context.Context, req *guilddto.SetGuildAnchorTypeReq) (*guilddto.SetGuildAnchorTypeRes, error) {
	if req == nil || req.GuildId == 0 || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.AnchorType != userentity.UserTypeAnchor && req.AnchorType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if guilddao.GetGuildById(req.GuildId) == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	user := userinfodao.GetUserInfoByUserId(req.UserId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if user.UserType != userentity.UserTypeAnchor && user.UserType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if liveroomdao.GetAnchorGuildId(req.UserId) != req.GuildId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if user.UserType == req.AnchorType {
		return &guilddto.SetGuildAnchorTypeRes{Success: true}, nil
	}
	user.SetUserType(req.AnchorType)
	liveroom.RefreshRoomListCache(ctx)
	return &guilddto.SetGuildAnchorTypeRes{Success: true}, nil
}
