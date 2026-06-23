package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/module/shortvideo"
)

const ShortVideoUrl = "/shortVideo"

type ShortVideoController struct{}

func initShortVideoController() {
	httpserver.RegCMS(ShortVideoUrl, &ShortVideoController{})
}

func (c *ShortVideoController) ShortVideoList(ctx context.Context, req *shortvideodto.ShortVideoListReq) (*httpserver.CMSQueryResp, error) {
	return shortvideo.GetShortVideoList(ctx, req)
}

func (c *ShortVideoController) CreateShortVideo(ctx context.Context, req *shortvideodto.CreateShortVideoReq) (*shortvideodto.CreateShortVideoRes, error) {
	return shortvideo.CreateShortVideo(ctx, req)
}

func (c *ShortVideoController) UpdateShortVideo(ctx context.Context, req *shortvideodto.UpdateShortVideoReq) (*shortvideodto.UpdateShortVideoRes, error) {
	return shortvideo.UpdateShortVideo(ctx, req)
}

func (c *ShortVideoController) DeleteShortVideo(ctx context.Context, req *shortvideodto.DeleteShortVideoReq) (*shortvideodto.DeleteShortVideoRes, error) {
	return shortvideo.DeleteShortVideo(ctx, req)
}

func (c *ShortVideoController) OnShelfShortVideo(ctx context.Context, req *shortvideodto.OnShelfShortVideoReq) (*shortvideodto.OnShelfShortVideoRes, error) {
	return shortvideo.OnShelfShortVideo(ctx, req)
}

func (c *ShortVideoController) OffShelfShortVideo(ctx context.Context, req *shortvideodto.OffShelfShortVideoReq) (*shortvideodto.OffShelfShortVideoRes, error) {
	return shortvideo.OffShelfShortVideo(ctx, req)
}

func (c *ShortVideoController) GetShortVideoCfg(ctx context.Context, req *shortvideodto.GetShortVideoCfgReq) (*shortvideodto.GetShortVideoCfgRes, error) {
	return shortvideo.GetShortVideoCfg(ctx, req)
}

func (c *ShortVideoController) SaveShortVideoCfg(ctx context.Context, req *shortvideodto.SaveShortVideoCfgReq) (*shortvideodto.SaveShortVideoCfgRes, error) {
	return shortvideo.SaveShortVideoCfg(ctx, req)
}

func (c *ShortVideoController) UploadShortVideo(ctx context.Context, req *shortvideodto.UploadShortVideoCMSReq) (*shortvideodto.UploadShortVideoCMSRes, error) {
	return shortvideo.UploadShortVideoCMS(ctx, req)
}

func (c *ShortVideoController) ShortVideoWatchList(ctx context.Context, req *shortvideodto.ShortVideoWatchListReq) (*httpserver.CMSQueryResp, error) {
	return shortvideo.GetShortVideoWatchList(ctx, req)
}

func (c *ShortVideoController) ShortVideoCategoryList(ctx context.Context, req *shortvideodto.ShortVideoCategoryListReq) (*httpserver.CMSQueryResp, error) {
	return shortvideo.GetShortVideoCategoryList(ctx, req)
}

func (c *ShortVideoController) CreateShortVideoCategory(ctx context.Context, req *shortvideodto.CreateShortVideoCategoryReq) (*shortvideodto.CreateShortVideoCategoryRes, error) {
	return shortvideo.CreateShortVideoCategory(ctx, req)
}

func (c *ShortVideoController) UpdateShortVideoCategory(ctx context.Context, req *shortvideodto.UpdateShortVideoCategoryReq) (*shortvideodto.UpdateShortVideoCategoryRes, error) {
	return shortvideo.UpdateShortVideoCategory(ctx, req)
}

func (c *ShortVideoController) DeleteShortVideoCategory(ctx context.Context, req *shortvideodto.DeleteShortVideoCategoryReq) (*shortvideodto.DeleteShortVideoCategoryRes, error) {
	return shortvideo.DeleteShortVideoCategory(ctx, req)
}
