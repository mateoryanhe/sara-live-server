package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/metapixeldto"
	"xr-game-server/module/metapixel"
)

const MetaPixelCMSUrl = "/metaPixel"

type MetaPixelCMSController struct{}

func initMetaPixelCMSController() {
	httpserver.RegCMS(MetaPixelCMSUrl, &MetaPixelCMSController{})
}

func (c *MetaPixelCMSController) GetMetaPixelCfg(ctx context.Context, req *metapixeldto.GetMetaPixelCfgReq) (*metapixeldto.GetMetaPixelCfgRes, error) {
	return metapixel.GetMetaPixelCfg(ctx, req)
}

func (c *MetaPixelCMSController) SaveMetaPixelCfg(ctx context.Context, req *metapixeldto.SaveMetaPixelCfgReq) (*metapixeldto.SaveMetaPixelCfgRes, error) {
	return metapixel.SaveMetaPixelCfg(ctx, req)
}
