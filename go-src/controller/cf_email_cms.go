package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/cfemaildto"
	"xr-game-server/module/auth"
)

const CfEmailCMSUrl = "/cfEmail"

type CfEmailCMSController struct{}

func initCfEmailCMSController() {
	httpserver.RegCMS(CfEmailCMSUrl, &CfEmailCMSController{})
}

func (c *CfEmailCMSController) GetCfEmailCfg(ctx context.Context, req *cfemaildto.GetCfEmailCfgReq) (*cfemaildto.GetCfEmailCfgRes, error) {
	return auth.GetCfEmailCfg(ctx, req)
}

func (c *CfEmailCMSController) SaveCfEmailCfg(ctx context.Context, req *cfemaildto.SaveCfEmailCfgReq) (*cfemaildto.SaveCfEmailCfgRes, error) {
	return auth.SaveCfEmailCfg(ctx, req)
}
