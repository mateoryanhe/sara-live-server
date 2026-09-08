package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/module/coinmerchant"
)

const CoinMerchantRechargeCfgCMSUrl = "/coinMerchantRechargeCfg"

type CoinMerchantRechargeCfgController struct{}

func initCoinMerchantRechargeCfgController() {
	httpserver.RegCMS(CoinMerchantRechargeCfgCMSUrl, &CoinMerchantRechargeCfgController{})
}

func (c *CoinMerchantRechargeCfgController) CoinMerchantRechargeCfgList(ctx context.Context, req *coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListReq) (*httpserver.CMSQueryResp, error) {
	return coinmerchant.ListCoinMerchantRechargeCfgs(ctx, req)
}

func (c *CoinMerchantRechargeCfgController) CreateCoinMerchantRechargeCfg(ctx context.Context, req *coinmerchantrechargecfgdto.CreateCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.CreateCoinMerchantRechargeCfgRes, error) {
	return coinmerchant.CreateCoinMerchantRechargeCfg(ctx, req)
}

func (c *CoinMerchantRechargeCfgController) UpdateCoinMerchantRechargeCfg(ctx context.Context, req *coinmerchantrechargecfgdto.UpdateCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.UpdateCoinMerchantRechargeCfgRes, error) {
	return coinmerchant.UpdateCoinMerchantRechargeCfg(ctx, req)
}

func (c *CoinMerchantRechargeCfgController) DeleteCoinMerchantRechargeCfg(ctx context.Context, req *coinmerchantrechargecfgdto.DeleteCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.DeleteCoinMerchantRechargeCfgRes, error) {
	return coinmerchant.DeleteCoinMerchantRechargeCfg(ctx, req)
}

func (c *CoinMerchantRechargeCfgController) OnShelfCoinMerchantRechargeCfg(ctx context.Context, req *coinmerchantrechargecfgdto.OnShelfCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.OnShelfCoinMerchantRechargeCfgRes, error) {
	return coinmerchant.OnShelfCoinMerchantRechargeCfg(ctx, req)
}

func (c *CoinMerchantRechargeCfgController) OffShelfCoinMerchantRechargeCfg(ctx context.Context, req *coinmerchantrechargecfgdto.OffShelfCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.OffShelfCoinMerchantRechargeCfgRes, error) {
	return coinmerchant.OffShelfCoinMerchantRechargeCfg(ctx, req)
}
