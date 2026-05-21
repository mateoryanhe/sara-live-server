package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/bannerdto"
	"xr-game-server/module/banner"
)

const BannerUrl = "/banner"

type BannerController struct{}

func initBannerController() {
	httpserver.RegCMS(BannerUrl, &BannerController{})
}

func (c *BannerController) BannerList(ctx context.Context, req *bannerdto.BannerListReq) (*httpserver.CMSQueryResp, error) {
	return banner.GetBannerList(ctx, req)
}

func (c *BannerController) CreateBanner(ctx context.Context, req *bannerdto.CreateBannerReq) (*bannerdto.CreateBannerRes, error) {
	return banner.CreateBanner(ctx, req)
}

func (c *BannerController) UpdateBanner(ctx context.Context, req *bannerdto.UpdateBannerReq) (*bannerdto.UpdateBannerRes, error) {
	return banner.UpdateBanner(ctx, req)
}

func (c *BannerController) DeleteBanner(ctx context.Context, req *bannerdto.DeleteBannerReq) (*bannerdto.DeleteBannerRes, error) {
	return banner.DeleteBanner(ctx, req)
}

func (c *BannerController) OnShelfBanner(ctx context.Context, req *bannerdto.OnShelfBannerReq) (*bannerdto.OnShelfBannerRes, error) {
	return banner.OnShelfBanner(ctx, req)
}

func (c *BannerController) OffShelfBanner(ctx context.Context, req *bannerdto.OffShelfBannerReq) (*bannerdto.OffShelfBannerRes, error) {
	return banner.OffShelfBanner(ctx, req)
}
