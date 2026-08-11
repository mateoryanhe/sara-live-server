package auth

import (
	"context"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
	"xr-game-server/module/accountcfg"
)

// CancelAccountByCode 通过注销码注销账号(官网公开接口)
func CancelAccountByCode(ctx context.Context, req *userinfodto.CancelAccountByCodeReq) (*userinfodto.CancelAccountByCodeRes, error) {
	if !accountcfg.IsCancelAccountByCodeEnabled() {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	cancelCode := strings.TrimSpace(req.CancelCode)
	if cancelCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := userinfodao.FindUserIdByCancelCode(cancelCode)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	ext := userinfodao.GetUserExtByUserId(userId)
	if strings.TrimSpace(ext.CancelCode) != cancelCode || !userinfodao.IsCancelCodeValid(ext) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	dbAcc := accountdao.GetAccountById(userId)
	if dbAcc == nil || dbAcc.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	account := accountdao.GetAccountFromCache(dbAcc.OpenId, dbAcc.Channel, userId)
	if account == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if account.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}

	ok, err := CancelUser(ctx, &accountdto.CancelReq{
		AccountId: userId,
		OpenId:    dbAcc.OpenId,
		Channel:   dbAcc.Channel,
	})
	if err != nil || !ok {
		return nil, err
	}
	return &userinfodto.CancelAccountByCodeRes{Success: true}, nil
}
