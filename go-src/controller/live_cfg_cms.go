package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/livecfgdto"
	"xr-game-server/module/livecfg"
)

const LiveCfgCMSUrl = "/liveCfg"

type LiveCfgCMSController struct{}

func initLiveCfgCMSController() {
	httpserver.RegCMS(LiveCfgCMSUrl, &LiveCfgCMSController{})
}

func (c *LiveCfgCMSController) GetLiveCfg(ctx context.Context, req *livecfgdto.GetLiveCfgReq) (*livecfgdto.GetLiveCfgRes, error) {
	return livecfg.GetLiveCfg(ctx, req)
}

func (c *LiveCfgCMSController) SaveLiveCfg(ctx context.Context, req *livecfgdto.SaveLiveCfgReq) (*livecfgdto.SaveLiveCfgRes, error) {
	return livecfg.SaveLiveCfg(ctx, req)
}
