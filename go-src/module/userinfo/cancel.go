package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

func CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	accountdao.CancelAccount(req.AccountId)
	return true, nil
}

func UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	accountdao.UnCancelAccount(req.AccountId)
	return true, nil
}

// CancelAccount App端销户(当前登录用户)
func CancelAccount(ctx context.Context, _ *userinfodto.CancelAccountReq) (*userinfodto.CancelAccountRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	accountdao.CancelAccount(userId)
	return &userinfodto.CancelAccountRes{Success: true}, nil
}
