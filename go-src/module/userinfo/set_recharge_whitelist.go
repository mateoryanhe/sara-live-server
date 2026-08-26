package userinfo

import (
	"context"

	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
)

// SetRechargeWhitelist CMS 设置用户充值白名单(白名单用户 App 创建订单后直接到账)
func SetRechargeWhitelist(_ context.Context, req *accountdto.SetRechargeWhitelistReq) (*accountdto.SetRechargeWhitelistRes, error) {
	ext := userinfodao.GetUserExtByUserId(req.AccountId)
	if ext.RechargeWhitelist == req.RechargeWhitelist {
		return &accountdto.SetRechargeWhitelistRes{Success: true}, nil
	}
	ext.SetRechargeWhitelist(req.RechargeWhitelist)
	userinfodao.PublishUserExt(ext)
	return &accountdto.SetRechargeWhitelistRes{Success: true}, nil
}
