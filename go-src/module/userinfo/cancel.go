package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

func CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	AddIdToCache(req.AccountId)
	db := accountdao.GetAccountById(req.AccountId)
	accountCache := accountdao.GetAccountBy(db.OpenId, db.Channel)
	accountCache.SetCancel(true)
	push.Kick(req.AccountId)
	return true, nil
}

func UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	AddIdToCache(req.AccountId)
	db := accountdao.GetAccountById(req.AccountId)
	accountCache := accountdao.GetAccountBy(db.OpenId, db.Channel)
	accountCache.SetCancel(false)
	return true, nil
}

// CancelAccount App端销户(当前登录用户)
func CancelAccount(ctx context.Context, _ *userinfodto.CancelAccountReq) (*userinfodto.CancelAccountRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	AddIdToCache(userId)
	db := accountdao.GetAccountById(userId)
	accountCache := accountdao.GetAccountBy(db.OpenId, db.Channel)
	accountCache.SetCancel(true)
	push.Kick(userId)
	return &userinfodto.CancelAccountRes{Success: true}, nil
}
