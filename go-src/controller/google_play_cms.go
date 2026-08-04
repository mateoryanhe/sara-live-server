package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/googleplaydto"
	"xr-game-server/module/recharge"
)

const GooglePlayCMSUrl = "/googlePlay"

type GooglePlayCMSController struct{}

func initGooglePlayCMSController() {
	httpserver.RegCMS(GooglePlayCMSUrl, &GooglePlayCMSController{})
}

func (c *GooglePlayCMSController) GetGooglePlayCfg(ctx context.Context, req *googleplaydto.GetGooglePlayCfgReq) (*googleplaydto.GetGooglePlayCfgRes, error) {
	return recharge.GetGooglePlayCfg(ctx, req)
}

func (c *GooglePlayCMSController) SaveGooglePlayCfg(ctx context.Context, req *googleplaydto.SaveGooglePlayCfgReq) (*googleplaydto.SaveGooglePlayCfgRes, error) {
	return recharge.SaveGooglePlayCfg(ctx, req)
}
