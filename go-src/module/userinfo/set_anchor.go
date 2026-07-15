package userinfo

import (
	"context"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// SetAnchor CMS 将用户设为主播(仅允许从未主播变为主播,不可回退)
func SetAnchor(_ context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	if user.UserType != entity.UserTypeNormal {
		return nil, errercode.CreateCode(errercode.UserAlreadyAnchor)
	}
	user.SetUserType(entity.UserTypeAnchor)
	liveroom.EnsureAnchorRoom(req.AccountId, user.GuildId)
	AddIdToCache(req.AccountId)
	return &accountdto.SetAnchorRes{Success: true}, nil
}
