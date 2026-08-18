package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liverevenuesharecfgdto"
	"xr-game-server/module/liverevenuesharecfg"
)

const LiveRevenueShareCfgCMSUrl = "/liveRevenueShareCfg"

type LiveRevenueShareCfgCMSController struct{}

func initLiveRevenueShareCfgCMSController() {
	httpserver.RegCMS(LiveRevenueShareCfgCMSUrl, &LiveRevenueShareCfgCMSController{})
}

func (c *LiveRevenueShareCfgCMSController) GetLiveRevenueShareCfg(ctx context.Context, req *liverevenuesharecfgdto.GetLiveRevenueShareCfgReq) (*liverevenuesharecfgdto.GetLiveRevenueShareCfgRes, error) {
	return liverevenuesharecfg.GetLiveRevenueShareCfg(ctx, req)
}

func (c *LiveRevenueShareCfgCMSController) SaveLiveRevenueShareCfg(ctx context.Context, req *liverevenuesharecfgdto.SaveLiveRevenueShareCfgReq) (*liverevenuesharecfgdto.SaveLiveRevenueShareCfgRes, error) {
	return liverevenuesharecfg.SaveLiveRevenueShareCfg(ctx, req)
}
