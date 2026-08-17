package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/simulatorcpukeyworddto"
	"xr-game-server/module/simulatorcpukeyword"
)

const SimulatorCpuKeywordCMSUrl = "/simulatorCpuKeyword"

type SimulatorCpuKeywordCMSController struct{}

func initSimulatorCpuKeywordCMSController() {
	httpserver.RegCMS(SimulatorCpuKeywordCMSUrl, &SimulatorCpuKeywordCMSController{})
}

func (c *SimulatorCpuKeywordCMSController) SimulatorCpuKeywordList(ctx context.Context, req *simulatorcpukeyworddto.SimulatorCpuKeywordListReq) (*httpserver.CMSQueryResp, error) {
	return simulatorcpukeyword.GetList(ctx, req)
}

func (c *SimulatorCpuKeywordCMSController) CreateSimulatorCpuKeyword(ctx context.Context, req *simulatorcpukeyworddto.CreateSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.CreateSimulatorCpuKeywordRes, error) {
	return simulatorcpukeyword.Create(ctx, req)
}

func (c *SimulatorCpuKeywordCMSController) UpdateSimulatorCpuKeyword(ctx context.Context, req *simulatorcpukeyworddto.UpdateSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.UpdateSimulatorCpuKeywordRes, error) {
	return simulatorcpukeyword.Update(ctx, req)
}

func (c *SimulatorCpuKeywordCMSController) DeleteSimulatorCpuKeyword(ctx context.Context, req *simulatorcpukeyworddto.DeleteSimulatorCpuKeywordReq) (*simulatorcpukeyworddto.DeleteSimulatorCpuKeywordRes, error) {
	return simulatorcpukeyword.Delete(ctx, req)
}
