package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/fiatcurrencydto"
	"xr-game-server/module/fiatcurrency"
)

const FiatCurrencyUrl = "/fiatCurrency"

type FiatCurrencyCMSController struct{}

func initFiatCurrencyCMSController() {
	httpserver.RegCMS(FiatCurrencyUrl, &FiatCurrencyCMSController{})
}

func (c *FiatCurrencyCMSController) FiatCurrencyList(ctx context.Context, req *fiatcurrencydto.FiatCurrencyListReq) (res *httpserver.CMSQueryResp, err error) {
	return fiatcurrency.GetList(ctx, req)
}

func (c *FiatCurrencyCMSController) CreateFiatCurrency(ctx context.Context, req *fiatcurrencydto.CreateFiatCurrencyReq) (res *fiatcurrencydto.CreateFiatCurrencyRes, err error) {
	return fiatcurrency.Create(ctx, req)
}

func (c *FiatCurrencyCMSController) UpdateFiatCurrency(ctx context.Context, req *fiatcurrencydto.UpdateFiatCurrencyReq) (res *fiatcurrencydto.UpdateFiatCurrencyRes, err error) {
	return fiatcurrency.Update(ctx, req)
}

func (c *FiatCurrencyCMSController) DeleteFiatCurrency(ctx context.Context, req *fiatcurrencydto.DeleteFiatCurrencyReq) (res *fiatcurrencydto.DeleteFiatCurrencyRes, err error) {
	return fiatcurrency.Delete(ctx, req)
}

func (c *FiatCurrencyCMSController) ReloadFiatCurrencyCache(ctx context.Context, req *fiatcurrencydto.ReloadFiatCurrencyCacheReq) (res *fiatcurrencydto.ReloadFiatCurrencyCacheRes, err error) {
	return fiatcurrency.ReloadCfgCache(ctx, req)
}

func (c *FiatCurrencyCMSController) ReloadFiatExchangeRateCache(ctx context.Context, req *fiatcurrencydto.ReloadFiatExchangeRateCacheReq) (res *fiatcurrencydto.ReloadFiatExchangeRateCacheRes, err error) {
	return fiatcurrency.ReloadExchangeRateCache(ctx, req)
}

func (c *FiatCurrencyCMSController) GetFiatExchangeRate(ctx context.Context, req *fiatcurrencydto.GetFiatExchangeRateReq) (res *fiatcurrencydto.GetFiatExchangeRateRes, err error) {
	return fiatcurrency.GetExchangeRate(ctx, req)
}
