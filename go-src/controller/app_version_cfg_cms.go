package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/appversioncfgdto"
	"xr-game-server/module/appversioncfg"
)

const AppVersionCfgCMSUrl = "/appVersionCfg"

type AppVersionCfgCMSController struct{}

func initAppVersionCfgCMSController() {
	httpserver.RegCMS(AppVersionCfgCMSUrl, &AppVersionCfgCMSController{})
}

func (c *AppVersionCfgCMSController) GetAppVersionCfg(ctx context.Context, req *appversioncfgdto.GetAppVersionCfgReq) (*appversioncfgdto.GetAppVersionCfgRes, error) {
	return appversioncfg.GetAppVersionCfg(ctx, req)
}

func (c *AppVersionCfgCMSController) SaveAppVersionCfg(ctx context.Context, req *appversioncfgdto.SaveAppVersionCfgReq) (*appversioncfgdto.SaveAppVersionCfgRes, error) {
	return appversioncfg.SaveAppVersionCfg(ctx, req)
}
