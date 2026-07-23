package userinfo

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/errercode"
)

// CancelAccountByCode 通过注销码注销账号(官网公开接口)
func CancelAccountByCode(ctx context.Context, req *userinfodto.CancelAccountByCodeReq) (*userinfodto.CancelAccountByCodeRes, error) {
	ip := g.RequestFromCtx(ctx).GetClientIp()
	if err := checkCancelByCodeIPLimit(ip); err != nil {
		return nil, err
	}

	cancelCode := strings.TrimSpace(req.CancelCode)
	if cancelCode == "" {
		return nil, markCancelByCodeFailure(ip, errercode.InvalidParam)
	}
	userId := userinfodao.FindUserIdByCancelCode(cancelCode)
	if userId == 0 {
		return nil, markCancelByCodeFailure(ip, errercode.InvalidParam)
	}
	account := accountdao.GetAccountById(userId)
	if account == nil {
		return nil, markCancelByCodeFailure(ip, errercode.InvalidParam)
	}
	if account.Cancel {
		return nil, markCancelByCodeFailure(ip, errercode.AccountCanceled)
	}
	ok, err := CancelUser(ctx, &accountdto.CancelReq{AccountId: userId})
	if err != nil || !ok {
		return nil, err
	}
	clearCancelByCodeFailure(ip)
	return &userinfodto.CancelAccountByCodeRes{Success: true}, nil
}
