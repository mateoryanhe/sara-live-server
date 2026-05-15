package account

import (
	"context"
	"time"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
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
