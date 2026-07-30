package liveroom

import (
	"context"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// SetAnchor CMS 将用户设为主播(仅允许从未主播变为主播,不可回退)
func SetAnchor(_ context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	return setUserAsAnchor(req.AccountId, entity.UserTypeAnchor)
}

// SetSeniorAnchor CMS 将用户设为高级主播(仅允许普通用户,不可回退)
func SetSeniorAnchor(_ context.Context, req *accountdto.SetSeniorAnchorReq) (*accountdto.SetSeniorAnchorRes, error) {
	_, err := setUserAsAnchor(req.AccountId, entity.UserTypeSeniorAnchor)
	if err != nil {
		return nil, err
	}
	return &accountdto.SetSeniorAnchorRes{Success: true}, nil
}

func setUserAsAnchor(accountId uint64, userType uint8) (*accountdto.SetAnchorRes, error) {
	user := userinfodao.GetUserInfoByUserId(accountId)
	if user.UserType != entity.UserTypeNormal {
		return nil, errercode.CreateCode(errercode.UserAlreadyAnchor)
	}
	user.SetUserType(userType)
	EnsureAnchorRoom(accountId, user.GuildId)
	RefreshRoomListCache(gctx.New())
	return &accountdto.SetAnchorRes{Success: true}, nil
}
