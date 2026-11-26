package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
)

const LiveRoomShareCMSUrl = "/liveRoomShare"

type LiveRoomShareCMSController struct{}

func initLiveRoomShareCMSController() {
	httpserver.RegCMS(LiveRoomShareCMSUrl, &LiveRoomShareCMSController{})
}

func (c *LiveRoomShareCMSController) GetLiveRoomShareCfg(ctx context.Context, req *shareappdownloaddto.GetLiveRoomShareCfgReq) (*shareappdownloaddto.GetLiveRoomShareCfgRes, error) {
	return shareappdownload.GetLiveRoomShareCfg(ctx, req)
}

func (c *LiveRoomShareCMSController) SaveLiveRoomShareCfg(ctx context.Context, req *shareappdownloaddto.SaveLiveRoomShareCfgReq) (*shareappdownloaddto.SaveLiveRoomShareCfgRes, error) {
	return shareappdownload.SaveLiveRoomShareCfg(ctx, req)
}