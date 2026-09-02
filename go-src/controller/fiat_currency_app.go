package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/fiatcurrencydto"
	"xr-game-server/module/fiatcurrency"
)

const FiatCurrencyAppUrl = "/fiatCurrency"

type FiatCurrencyAppController struct{}

func initFiatCurrencyAppController() {
	httpserver.RegNonAuthAPI(FiatCurrencyAppUrl, &FiatCurrencyAppController{})
}

func (c *FiatCurrencyAppController) FiatCurrencyListForApp(ctx context.Context, req *fiatcurrencydto.AppFiatCurrencyListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	return fiatcurrency.GetAppList(ctx, req)
}
