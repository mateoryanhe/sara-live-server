package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/appsflyercfgdto"
	"xr-game-server/module/appsflyercfg"
)

const AppsFlyerCMSUrl = "/appsFlyer"

type AppsFlyerCMSController struct{}

func initAppsFlyerCMSController() {
	httpserver.RegCMS(AppsFlyerCMSUrl, &AppsFlyerCMSController{})
}

func (c *AppsFlyerCMSController) GetAppsFlyerCfg(ctx context.Context, req *appsflyercfgdto.GetAppsFlyerCfgReq) (*appsflyercfgdto.GetAppsFlyerCfgRes, error) {
	return appsflyercfg.GetAppsFlyerCfg(ctx, req)
}

func (c *AppsFlyerCMSController) SaveAppsFlyerCfg(ctx context.Context, req *appsflyercfgdto.SaveAppsFlyerCfgReq) (*appsflyercfgdto.SaveAppsFlyerCfgRes, error) {
	return appsflyercfg.SaveAppsFlyerCfg(ctx, req)
}