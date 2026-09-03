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

// RechargeCfgListByUserId App按用户ID查询充值配置列表(已上架,无需鉴权)
func (c *RechargeCfgAppPublicController) RechargeCfgListByUserId(ctx context.Context, req *rechargecfgdto.AppRechargeCfgListByUserIdReq) (res *rechargecfgdto.AppRechargeCfgListRes, err error) {
	return recharge.GetAppListByUserId(ctx, req)
}
