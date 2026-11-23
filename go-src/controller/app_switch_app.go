package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/appswitchdto"
	"xr-game-server/module/appswitch"
)

const AppSwitchAppUrl = "/appSwitch"

type AppSwitchAppController struct{}

func initAppSwitchAppController() {
	httpserver.RegAPI(AppSwitchAppUrl, &AppSwitchAppController{})
}

func (c *AppSwitchAppController) GetAppSwitches(ctx context.Context, req *appswitchdto.GetAppSwitchesReq) (*appswitchdto.GetAppSwitchesRes, error) {
	return appswitch.GetAppSwitches(ctx, req)
}