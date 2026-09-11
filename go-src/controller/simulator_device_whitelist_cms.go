package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/simulatordevicewhitelistdto"
	"xr-game-server/module/simulatordevicewhitelist"
)

const SimulatorDeviceWhitelistCMSUrl = "/simulatorDeviceWhitelist"

type SimulatorDeviceWhitelistCMSController struct{}

func initSimulatorDeviceWhitelistCMSController() {
	httpserver.RegCMS(SimulatorDeviceWhitelistCMSUrl, &SimulatorDeviceWhitelistCMSController{})
}

func (c *SimulatorDeviceWhitelistCMSController) SimulatorDeviceWhitelistList(ctx context.Context, req *simulatordevicewhitelistdto.SimulatorDeviceWhitelistListReq) (*httpserver.CMSQueryResp, error) {
	return simulatordevicewhitelist.GetList(ctx, req)
}

func (c *SimulatorDeviceWhitelistCMSController) CreateSimulatorDeviceWhitelist(ctx context.Context, req *simulatordevicewhitelistdto.CreateSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.CreateSimulatorDeviceWhitelistRes, error) {
	return simulatordevicewhitelist.Create(ctx, req)
}

func (c *SimulatorDeviceWhitelistCMSController) UpdateSimulatorDeviceWhitelist(ctx context.Context, req *simulatordevicewhitelistdto.UpdateSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.UpdateSimulatorDeviceWhitelistRes, error) {
	return simulatordevicewhitelist.Update(ctx, req)
}

func (c *SimulatorDeviceWhitelistCMSController) DeleteSimulatorDeviceWhitelist(ctx context.Context, req *simulatordevicewhitelistdto.DeleteSimulatorDeviceWhitelistReq) (*simulatordevicewhitelistdto.DeleteSimulatorDeviceWhitelistRes, error) {
	return simulatordevicewhitelist.Delete(ctx, req)
}
