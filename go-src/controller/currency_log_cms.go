package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/currencylogdto"
	"xr-game-server/module/currencylog"
)

const (
	CurrencyLogUrl = "/currencyLog"
)

type CurrencyLogController struct{}

func initCurrencyLogController() {
	httpserver.RegCMS(CurrencyLogUrl, &CurrencyLogController{})
}

// CMSCurrencyLogList CMS分页查询货币流水
func (c *CurrencyLogController) CMSCurrencyLogList(ctx context.Context, req *currencylogdto.CMSCurrencyLogListReq) (res *httpserver.CMSQueryResp, err error) {
	return currencylog.GetCMSList(ctx, req)
}
