package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/module/shortvideo"
)

const ShortVideoAppUrl = "/shortVideo"

type ShortVideoAppController struct{}

func initShortVideoAppController() {
	httpserver.RegAPI(ShortVideoAppUrl, &ShortVideoAppController{})
}

func (c *ShortVideoAppController) AppShortVideoList(ctx context.Context, req *shortvideodto.AppShortVideoListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return shortvideo.GetAppShortVideoList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoViewList(ctx context.Context, req *shortvideodto.AppShortVideoViewListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return shortvideo.GetAppShortVideoViewList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoPublishList(ctx context.Context, req *shortvideodto.AppShortVideoPublishListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return shortvideo.GetAppShortVideoPublishList(ctx, req)
}

func (c *ShortVideoAppController) LikeShortVideo(ctx context.Context, req *shortvideodto.LikeShortVideoReq) (*shortvideodto.LikeShortVideoRes, error) {
	return shortvideo.LikeShortVideo(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoCfg(ctx context.Context, req *shortvideodto.AppShortVideoCfgReq) (*shortvideodto.AppShortVideoCfgRes, error) {
	return shortvideo.GetAppShortVideoCfg(ctx, req)
}

func (c *ShortVideoAppController) WatchBillShortVideo(ctx context.Context, req *shortvideodto.WatchBillShortVideoReq) (*shortvideodto.WatchBillShortVideoRes, error) {
	return shortvideo.WatchBillShortVideo(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoWatchList(ctx context.Context, req *shortvideodto.AppShortVideoWatchListReq) (*shortvideodto.AppShortVideoWatchListRes, error) {
	return shortvideo.GetAppShortVideoWatchList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoCategoryList(ctx context.Context, req *shortvideodto.AppShortVideoCategoryListReq) (*shortvideodto.AppShortVideoCategoryListRes, error) {
	return shortvideo.GetAppShortVideoCategoryList(ctx, req)
}

func (c *ShortVideoAppController) AppPublishShortVideo(ctx context.Context, req *shortvideodto.AppPublishShortVideoReq) (*shortvideodto.AppPublishShortVideoRes, error) {
	return shortvideo.PublishShortVideoApp(ctx, req)
}
