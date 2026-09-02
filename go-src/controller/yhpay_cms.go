package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/yhpaydto"
	"xr-game-server/module/recharge"
)

const YhPayCMSUrl = "/yhpay"

type YhPayCMSController struct{}

func initYhPayCMSController() {
	httpserver.RegCMS(YhPayCMSUrl, &YhPayCMSController{})
}

func (c *YhPayCMSController) GetYhPayCfg(ctx context.Context, req *yhpaydto.GetYhPayCfgReq) (*yhpaydto.GetYhPayCfgRes, error) {
	return recharge.GetYhPayCfg(ctx, req)
}

func (c *YhPayCMSController) SaveYhPayCfg(ctx context.Context, req *yhpaydto.SaveYhPayCfgReq) (*yhpaydto.SaveYhPayCfgRes, error) {
	return recharge.SaveYhPayCfg(ctx, req)
}
