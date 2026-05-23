package userinfo

import (
	"context"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
)

// SetAnchor CMS 将用户设为主播(仅允许从未主播变为主播,不可回退)
func SetAnchor(_ context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	if accountdao.GetAccountById(req.AccountId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	user := userinfodao.GetUserInfoByUserId(req.AccountId)
	if user.IsAnchor {
		return nil, errercode.CreateCode(errercode.UserAlreadyAnchor)
	}
	user.SetIsAnchor(true)
	return &accountdto.SetAnchorRes{Success: true}, nil
}
