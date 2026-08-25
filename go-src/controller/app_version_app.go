package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/appversioncfgdto"
	"xr-game-server/module/appversioncfg"
)

const AppVersionAppUrl = "/appVersion"

type AppVersionAppController struct{}

func initAppVersionAppController() {
	httpserver.RegNonAuthAPI(AppVersionAppUrl, &AppVersionAppController{})
}

func (c *AppVersionAppController) AppVersionQuery(ctx context.Context, req *appversioncfgdto.AppVersionQueryReq) (*appversioncfgdto.AppVersionQueryRes, error) {
	return appversioncfg.AppVersionQuery(ctx, req)
}
