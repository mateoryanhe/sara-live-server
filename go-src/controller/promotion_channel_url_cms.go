package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/module/officialsitedeploy"
	"xr-game-server/module/promotionchannelurl"
)

const PromotionChannelUrlCMSUrl = "/promotionChannelUrl"

type PromotionChannelUrlCMSController struct{}

func initPromotionChannelUrlCMSController() {
	httpserver.RegCMS(PromotionChannelUrlCMSUrl, &PromotionChannelUrlCMSController{})
	httpserver.RegCMSHandler(PromotionChannelUrlCMSUrl, "/uploadPromotionChannelUrlApkChunk", handlePromotionChannelUrlUploadApkChunk)
}

func handlePromotionChannelUrlUploadApkChunk(r *ghttp.Request) {
	res, err := officialsitedeploy.UploadPromotionChannelUrlApkChunkFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *PromotionChannelUrlCMSController) PromotionChannelUrlList(ctx context.Context, req *promotionchannelurldto.PromotionChannelUrlListReq) (*promotionchannelurldto.PromotionChannelUrlListRes, error) {
	return promotionchannelurl.GetList(ctx, req)
}

func (c *PromotionChannelUrlCMSController) CreatePromotionChannelUrl(ctx context.Context, req *promotionchannelurldto.CreatePromotionChannelUrlReq) (*promotionchannelurldto.CreatePromotionChannelUrlRes, error) {
	return promotionchannelurl.Create(ctx, req)
}

func (c *PromotionChannelUrlCMSController) UpdatePromotionChannelUrl(ctx context.Context, req *promotionchannelurldto.UpdatePromotionChannelUrlReq) (*promotionchannelurldto.UpdatePromotionChannelUrlRes, error) {
	return promotionchannelurl.Update(ctx, req)
}

func (c *PromotionChannelUrlCMSController) DeletePromotionChannelUrl(ctx context.Context, req *promotionchannelurldto.DeletePromotionChannelUrlReq) (*promotionchannelurldto.DeletePromotionChannelUrlRes, error) {
	return promotionchannelurl.Delete(ctx, req)
}

func (c *PromotionChannelUrlCMSController) InitPromotionChannelUrlApkUpload(ctx context.Context, req *promotionchannelurldto.InitPromotionChannelUrlApkUploadReq) (*promotionchannelurldto.InitPromotionChannelUrlApkUploadRes, error) {
	return officialsitedeploy.InitPromotionChannelUrlApkUpload(ctx, req)
}

func (c *PromotionChannelUrlCMSController) CompletePromotionChannelUrlApkUpload(ctx context.Context, req *promotionchannelurldto.CompletePromotionChannelUrlApkUploadReq) (*promotionchannelurldto.CompletePromotionChannelUrlApkUploadRes, error) {
	return officialsitedeploy.CompletePromotionChannelUrlApkUpload(ctx, req)
}

func (c *PromotionChannelUrlCMSController) AbortPromotionChannelUrlApkUpload(ctx context.Context, req *promotionchannelurldto.AbortPromotionChannelUrlApkUploadReq) (*promotionchannelurldto.AbortPromotionChannelUrlApkUploadRes, error) {
	return officialsitedeploy.AbortPromotionChannelUrlApkUpload(ctx, req)
}

func (c *PromotionChannelUrlCMSController) PromotionChannelDailyStatList(ctx context.Context, req *promotionchannelurldto.PromotionChannelDailyStatListReq) (*promotionchannelurldto.PromotionChannelDailyStatListRes, error) {
	return promotionchannelurl.GetDailyStatList(ctx, req)
}

func (c *PromotionChannelUrlCMSController) H5AppPromotionDailyStatList(ctx context.Context, req *promotionchannelurldto.H5AppPromotionDailyStatListReq) (*promotionchannelurldto.H5AppPromotionDailyStatListRes, error) {
	return promotionchannelurl.GetH5AppPromotionDailyStatList(ctx, req)
}
