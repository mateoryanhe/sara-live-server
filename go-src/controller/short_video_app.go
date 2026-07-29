package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/module/shortvideo"
)

const ShortVideoAppUrl = "/shortVideo"

type ShortVideoAppController struct{}

func initShortVideoAppController() {
	httpserver.RegAPIHandler(ShortVideoAppUrl, "/appPublishShortVideo", handleAppPublishShortVideo)
	httpserver.RegAPI(ShortVideoAppUrl, &ShortVideoAppController{})
}

func handleAppPublishShortVideo(r *ghttp.Request) {
	res, err := shortvideo.PublishShortVideoAppFromRequest(r.Context(), r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *ShortVideoAppController) AppShortVideoList(ctx context.Context, req *shortvideodto.AppShortVideoListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return shortvideo.GetAppShortVideoList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoScroll(ctx context.Context, req *shortvideodto.AppShortVideoScrollReq) (*shortvideodto.AppShortVideoScrollRes, error) {
	return shortvideo.GetAppShortVideoScrollList(ctx, req)
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

func (c *ShortVideoAppController) WatchShortVideoStart(ctx context.Context, req *shortvideodto.WatchShortVideoStartReq) (*shortvideodto.WatchShortVideoStartRes, error) {
	return shortvideo.WatchShortVideoStart(ctx, req)
}

func (c *ShortVideoAppController) WatchShortVideoEnd(ctx context.Context, req *shortvideodto.WatchShortVideoEndReq) (*shortvideodto.WatchShortVideoEndRes, error) {
	return shortvideo.WatchShortVideoEnd(ctx, req)
}

func (c *ShortVideoAppController) PayShortVideo(ctx context.Context, req *shortvideodto.PayShortVideoReq) (*shortvideodto.PayShortVideoRes, error) {
	return shortvideo.PayShortVideo(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoWatchList(ctx context.Context, req *shortvideodto.AppShortVideoWatchListReq) (*shortvideodto.AppShortVideoListRes, error) {
	return shortvideo.GetAppShortVideoWatchList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoCategoryList(ctx context.Context, req *shortvideodto.AppShortVideoCategoryListReq) (*shortvideodto.AppShortVideoCategoryListRes, error) {
	return shortvideo.GetAppShortVideoCategoryList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoUploadRecordList(ctx context.Context, req *shortvideodto.AppShortVideoUploadRecordListReq) (*shortvideodto.AppShortVideoUploadRecordListRes, error) {
	return shortvideo.GetAppShortVideoUploadRecordList(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoPendingReviewList(ctx context.Context, req *shortvideodto.AppShortVideoPendingReviewListReq) (*shortvideodto.AppShortVideoUploadRecordListRes, error) {
	return shortvideo.GetAppShortVideoPendingReviewList(ctx, req)
}

func (c *ShortVideoAppController) AppDeleteShortVideo(ctx context.Context, req *shortvideodto.AppDeleteShortVideoReq) (*shortvideodto.AppDeleteShortVideoRes, error) {
	return shortvideo.DeleteShortVideoApp(ctx, req)
}

func (c *ShortVideoAppController) AppShortVideoStatList(ctx context.Context, req *shortvideodto.AppShortVideoStatListReq) (*shortvideodto.AppShortVideoStatListRes, error) {
	return shortvideo.GetAppShortVideoStatList(ctx, req)
}
