package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/haipaydto"
	"xr-game-server/module/recharge"
)

const HaiPayCMSUrl = "/haipay"

type HaiPayCMSController struct{}

func initHaiPayCMSController() {
	httpserver.RegCMS(HaiPayCMSUrl, &HaiPayCMSController{})
}

func (c *HaiPayCMSController) GetHaiPayCfg(ctx context.Context, req *haipaydto.GetHaiPayCfgReq) (*haipaydto.GetHaiPayCfgRes, error) {
	return recharge.GetHaiPayCfg(ctx, req)
}

func (c *HaiPayCMSController) SaveHaiPayCfg(ctx context.Context, req *haipaydto.SaveHaiPayCfgReq) (*haipaydto.SaveHaiPayCfgRes, error) {
	return recharge.SaveHaiPayCfg(ctx, req)
}
