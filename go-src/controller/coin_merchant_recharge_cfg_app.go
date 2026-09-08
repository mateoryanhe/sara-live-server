package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/module/coinmerchant"
)

const CoinMerchantRechargeCfgAppUrl = "/coinMerchantRechargeCfg"

type CoinMerchantRechargeCfgAppController struct{}

func initCoinMerchantRechargeCfgAppController() {
	httpserver.RegAPI(CoinMerchantRechargeCfgAppUrl, &CoinMerchantRechargeCfgAppController{})
}

// CoinMerchantRechargeCfgList App端查询币商充值档位列表(仅已上架)
func (c *CoinMerchantRechargeCfgAppController) CoinMerchantRechargeCfgList(ctx context.Context, req *coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgListReq) (*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgListRes, error) {
	return coinmerchant.GetAppList(ctx, req)
}
