package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/module/rechargecfg"
)

const (
	RechargeCfgAppUrl = "/rechargeCfg"
)

type RechargeCfgAppController struct{}

func initRechargeCfgAppController() {
	httpserver.RegAPI(RechargeCfgAppUrl, &RechargeCfgAppController{})
}

// RechargeCfgList App端查询充值配置列表(仅已上架)
func (c *RechargeCfgAppController) RechargeCfgList(ctx context.Context, req *rechargecfgdto.AppRechargeCfgListReq) (res *rechargecfgdto.AppRechargeCfgListRes, err error) {
	return rechargecfg.GetAppList(ctx, req)
}
