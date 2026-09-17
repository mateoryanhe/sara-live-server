package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/dto/fiatcurrencydto"
	"xr-game-server/module/coinmerchant"
	"xr-game-server/module/fiatcurrency"
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

// PaymentRegionList App端查询币商专用支付区域列表。
func (c *CoinMerchantRechargeCfgAppController) PaymentRegionList(ctx context.Context, req *fiatcurrencydto.CoinMerchantPaymentRegionListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	return fiatcurrency.GetCoinMerchantAppList(ctx, req)
}
