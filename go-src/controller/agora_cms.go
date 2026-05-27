package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/agoradto"
	"xr-game-server/module/agora"
)

const AgoraCMSUrl = "/agora"

type AgoraCMSController struct{}

func initAgoraCMSController() {
	httpserver.RegCMS(AgoraCMSUrl, &AgoraCMSController{})
}

func (c *AgoraCMSController) GetAgoraCfg(ctx context.Context, req *agoradto.GetAgoraCfgReq) (*agoradto.GetAgoraCfgRes, error) {
	return agora.GetAgoraCfg(ctx, req)
}

func (c *AgoraCMSController) SaveAgoraCfg(ctx context.Context, req *agoradto.SaveAgoraCfgReq) (*agoradto.SaveAgoraCfgRes, error) {
	return agora.SaveAgoraCfg(ctx, req)
}
