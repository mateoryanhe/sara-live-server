package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	account := accountdao.GetAccountFromCache(req.OpenId, req.Channel, req.AccountId)
	if account == nil {
		return false, errercode.CreateCode(errercode.InvalidParam)
	}
	account.SetCancel(true)
	invalidateAppToken(req.AccountId)
	push.Kick(req.AccountId)
	return true, nil
}

func UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	account := accountdao.GetAccountFromCache(req.OpenId, req.Channel, req.AccountId)
	if account == nil {
		return false, errercode.CreateCode(errercode.InvalidParam)
	}
	if !account.Cancel {
		return true, nil
	}
	if active := accountdao.FindActiveAccount(req.OpenId, req.Channel); active != nil && active.ID != req.AccountId {
		return false, errercode.CreateCode(errercode.AccountAlreadyExists)
	}
	account.SetCancel(false)
	return true, nil
}

// CancelAccount App端销户(需登录,无需请求体)
func CancelAccount(ctx context.Context, _ *userinfodto.CancelAccountReq) (*userinfodto.CancelAccountRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := doAppCancelAccount(userId); err != nil {
		return nil, err
	}
	return &userinfodto.CancelAccountRes{Success: true}, nil
}

func doAppCancelAccount(accountId uint64) error {
	if accountId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	dbAcc := accountdao.GetAccountById(accountId)
	if dbAcc == nil || dbAcc.ID == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	account := accountdao.GetAccountFromCache(dbAcc.OpenId, dbAcc.Channel, accountId)
	if account == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if account.Cancel {
		return errercode.CreateCode(errercode.AccountCanceled)
	}
	if err := checkAppCancelAccountGuard(dbAcc.OpenId, dbAcc.Channel); err != nil {
		return err
	}
	account.SetCancel(true)
	recordAppCancelAccountSuccess(dbAcc.OpenId, dbAcc.Channel)
	invalidateAppToken(accountId)
	push.Kick(accountId)
	return nil
}

func invalidateAppToken(userId uint64) {
	if userId == 0 {
		return
	}
	xrtoken.InitAppToken(userId, guid.S(), time.Now().Add(30*time.Minute))
	expireTime := time.Now().Add(-24 * 100 * time.Hour)
	entity.NewAppToken(userId, fmt.Sprintf("%v", userId), expireTime)
}
