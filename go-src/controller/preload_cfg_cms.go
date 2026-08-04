package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/preloadcfgdto"
	"xr-game-server/module/preload"
)

const PreloadCfgCMSUrl = "/preloadCfg"

type PreloadCfgCMSController struct{}

func initPreloadCfgCMSController() {
	httpserver.RegCMS(PreloadCfgCMSUrl, &PreloadCfgCMSController{})
}

func (c *PreloadCfgCMSController) GetPreloadCfg(ctx context.Context, req *preloadcfgdto.GetPreloadCfgReq) (*preloadcfgdto.GetPreloadCfgRes, error) {
	return preload.GetPreloadCfg(ctx, req)
}

func (c *PreloadCfgCMSController) SavePreloadCfg(ctx context.Context, req *preloadcfgdto.SavePreloadCfgReq) (*preloadcfgdto.SavePreloadCfgRes, error) {
	return preload.SavePreloadCfg(ctx, req)
}
