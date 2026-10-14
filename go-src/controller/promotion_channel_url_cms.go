package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/module/promotionchannelurl"
)

const PromotionChannelUrlCMSUrl = "/promotionChannelUrl"

type PromotionChannelUrlCMSController struct{}

func initPromotionChannelUrlCMSController() {
	httpserver.RegCMS(PromotionChannelUrlCMSUrl, &PromotionChannelUrlCMSController{})
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

func (c *PromotionChannelUrlCMSController) PromotionChannelDailyStatList(ctx context.Context, req *promotionchannelurldto.PromotionChannelDailyStatListReq) (*promotionchannelurldto.PromotionChannelDailyStatListRes, error) {
	return promotionchannelurl.GetDailyStatList(ctx, req)
}
