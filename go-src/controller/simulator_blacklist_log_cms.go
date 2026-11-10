package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/simulatorblacklistlogdto"
	"xr-game-server/module/simulatorblacklistlog"
)

const SimulatorBlacklistLogCMSUrl = "/simulatorBlacklistLog"

type SimulatorBlacklistLogCMSController struct{}

func initSimulatorBlacklistLogCMSController() {
	httpserver.RegCMS(SimulatorBlacklistLogCMSUrl, &SimulatorBlacklistLogCMSController{})
}

func (c *SimulatorBlacklistLogCMSController) ListSimulatorBlacklistLog(ctx context.Context, req *simulatorblacklistlogdto.CMSListSimulatorBlacklistLogReq) (*httpserver.CMSQueryResp, error) {
	return simulatorblacklistlog.GetCMSList(ctx, req)
}
