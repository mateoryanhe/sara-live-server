package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/bannerdto"
	"xr-game-server/module/banner"
)

const BannerAppUrl = "/banner"

type BannerAppController struct{}

func initBannerAppController() {
	httpserver.RegAPI(BannerAppUrl, &BannerAppController{})
}

func (c *BannerAppController) AppBannerList(ctx context.Context, req *bannerdto.AppBannerListReq) (*bannerdto.AppBannerListRes, error) {
	return banner.GetAppBannerList(ctx, req)
}
