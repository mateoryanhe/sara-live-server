package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/module/recharge"
)

// RechargeCfgAppPublicController App充值配置公开接口(无需鉴权)
type RechargeCfgAppPublicController struct{}

func initRechargeCfgAppPublicController() {
	httpserver.RegNonAuthAPI(RechargeCfgAppUrl, &RechargeCfgAppPublicController{})
}

// RechargeCfgListByCurrency App按币种查询充值配置列表(价格按实时汇率+加点换算,无需鉴权)
func (c *RechargeCfgAppPublicController) RechargeCfgListByCurrency(ctx context.Context, req *rechargecfgdto.AppRechargeCfgListByCurrencyReq) (res *rechargecfgdto.AppRechargeCfgListRes, err error) {
	return recharge.GetAppListByCurrency(ctx, req)
}
