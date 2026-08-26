package auth

import (
	"context"
	"time"

	"xr-game-server/core/push"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/errercode"
)

func Ban(ctx context.Context, req *accountdto.BanReq) (resp *accountdto.BanRes, e error) {
	account := accountdao.GetAccountFromCache(req.OpenId, req.Channel, req.AccountId)
	if account == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	account.SetBan(true)
	account.SetUpdatedAt(time.Now())
	now := time.Now()
	account.SetBanTime(&now)
	account.SetBanApplyTime(req.BanApplyTime)
	accountdao.PublishAccountList(req.OpenId, req.Channel)
	push.Kick(req.AccountId)
	invalidateAppToken(req.AccountId)

	return &accountdto.BanRes{}, nil
}

func UnBan(ctx context.Context, req *accountdto.UnBanReq) (bool, error) {
	account := accountdao.GetAccountFromCache(req.OpenId, req.Channel, req.AccountId)
	if account == nil {
		return false, errercode.CreateCode(errercode.InvalidParam)
	}
	account.SetBan(false)
	account.SetBanTime(nil)
	account.SetBanApplyTime(nil)
	account.SetUpdatedAt(time.Now())
	accountdao.PublishAccountList(req.OpenId, req.Channel)
	return true, nil
}
