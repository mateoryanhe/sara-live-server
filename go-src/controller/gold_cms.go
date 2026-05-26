package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/golddto"
	"xr-game-server/module/userinfo"
)

const GoldUrl = "/gold"

type GoldController struct{}

func initGoldController() {
	httpserver.RegCMS(GoldUrl, &GoldController{})
}

func (c *GoldController) Add(ctx context.Context, req *golddto.CMSAddGoldReq) (*golddto.CMSAddGoldRes, error) {
	return userinfo.GoldCMSAdd(ctx, req)
}

func (c *GoldController) Sub(ctx context.Context, req *golddto.CMSSubGoldReq) (*golddto.CMSSubGoldRes, error) {
	return userinfo.GoldCMSSub(ctx, req)
}
