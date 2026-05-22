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

func (c *ShortVideoAppController) LikeShortVideo(ctx context.Context, req *shortvideodto.LikeShortVideoReq) (*shortvideodto.LikeShortVideoRes, error) {
	return shortvideo.LikeShortVideo(ctx, req)
}
