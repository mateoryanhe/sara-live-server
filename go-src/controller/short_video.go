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
