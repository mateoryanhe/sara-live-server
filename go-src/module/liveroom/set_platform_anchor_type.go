package liveroom

import (
	"context"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// SetPlatformAnchorType CMS设置平台主播类型(仅支持普通/高级主播互切,且 guild_id=0)
func SetPlatformAnchorType(ctx context.Context, req *accountdto.SetPlatformAnchorTypeReq) (*accountdto.SetPlatformAnchorTypeRes, error) {
	if req == nil || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.AnchorType != userentity.UserTypeAnchor && req.AnchorType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user := userinfodao.GetUserInfoByUserId(req.UserId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if user.UserType != userentity.UserTypeAnchor && user.UserType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if liveroomdao.GetAnchorGuildId(req.UserId) != 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if user.UserType == req.AnchorType {
		return &accountdto.SetPlatformAnchorTypeRes{Success: true}, nil
	}
	user.SetUserType(req.AnchorType)
	userinfodao.PublishUserInfo(user)
	RefreshRoomListCache(ctx)
	return &accountdto.SetPlatformAnchorTypeRes{Success: true}, nil
}
