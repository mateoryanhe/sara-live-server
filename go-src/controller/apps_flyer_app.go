package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/appsflyercfgdto"
	"xr-game-server/module/appsflyercfg"
)

type AppsFlyerAppController struct{}

func initAppsFlyerAppController() {
	httpserver.RegNonAuthAPI(AppsFlyerCMSUrl, &AppsFlyerAppController{})
}

func (c *AppsFlyerAppController) GetCfgForApp(ctx context.Context, req *appsflyercfgdto.GetAppsFlyerCfgForAppReq) (*appsflyercfgdto.GetAppsFlyerCfgForAppRes, error) {
	return appsflyercfg.GetAppsFlyerCfgForApp(ctx, req)
}