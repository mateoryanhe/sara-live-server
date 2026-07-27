package auth

import (
	"context"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

// CancelAccountByCode 通过注销码注销账号(官网公开接口)
func CancelAccountByCode(ctx context.Context, req *userinfodto.CancelAccountByCodeReq) (*userinfodto.CancelAccountByCodeRes, error) {
	cancelCode := strings.TrimSpace(req.CancelCode)
	if cancelCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := userinfodao.FindUserIdByCancelCode(cancelCode)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	dbAcc := accountdao.GetAccountById(userId)
	if dbAcc == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
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
