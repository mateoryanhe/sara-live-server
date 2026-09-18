package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/golddto"
	"xr-game-server/module/userinfo"
	"xr-game-server/module/wallet"
)

const GoldUrl = "/gold"

type GoldController struct{}

func initGoldController() {
	httpserver.RegCMS(GoldUrl, &GoldController{})
}

func (c *GoldController) Add(ctx context.Context, req *golddto.CMSAddGoldReq) (*golddto.CMSAddGoldRes, error) {
	return userinfo.GoldCMSAdd(ctx, req)
}

func (c *GoldController) Sub(ctx context.Context, req *golddto.CMSSubGoldReq) (*golddto.CMSSubGoldRes, error) {
	return userinfo.GoldCMSSub(ctx, req)
}

// CMSCoinMerchantTransferLogList CMS 分页查询币商金币转账流水。
func (c *GoldController) CMSCoinMerchantTransferLogList(ctx context.Context, req *golddto.CMSCoinMerchantTransferLogListReq) (*httpserver.CMSQueryResp, error) {
	return wallet.GetCMSCoinMerchantTransferLogList(ctx, req)
}
