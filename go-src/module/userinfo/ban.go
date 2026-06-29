package userinfo

import (
	"context"
	"fmt"
	"time"
	"xr-game-server/core/push"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
)

func Ban(ctx context.Context, req *accountdto.BanReq) (resp *accountdto.BanRes, e error) {
	//默认永远封掉号
	db := accountdao.GetAccountById(req.AccountId)
	data := accountdao.GetAccountBy(db.OpenId, db.Channel)
	data.SetBan(true)
	data.SetUpdatedAt(time.Now())
	now := time.Now()
	data.SetBanTime(&now)
	data.SetBanApplyTime(req.BanApplyTime)
	AddIdToCache(req.AccountId)
	push.Kick(req.AccountId)
	xrtoken.DelToken(req.AccountId)
	expireTime := time.Now().Add(-24 * 100 * time.Hour)
	entity.NewAppToken(req.AccountId, fmt.Sprintf("%v", req.AccountId), expireTime)
	//踢掉用户
	//roles := roledao.GetRoleBy(req.AccountId)
	//for _, role := range roles {
	//	push.Kick(role.ID)
	//}

	return &accountdto.BanRes{}, nil
}
func UnBan(ctx context.Context, req *accountdto.UnBanReq) (bool, error) {
	db := accountdao.GetAccountById(req.AccountId)
	data := accountdao.GetAccountBy(db.OpenId, db.Channel)
	data.SetBan(false)
	data.SetBanTime(nil)
	data.SetBanApplyTime(nil)
	data.SetUpdatedAt(time.Now())
	return true, nil
}
