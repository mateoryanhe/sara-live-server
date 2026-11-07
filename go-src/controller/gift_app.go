package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/giftdto"
	"xr-game-server/module/gift"
)

const (
	GiftAppUrl = "/gift"
)

type GiftAppController struct{}

func initGiftAppController() {
	httpserver.RegAPI(GiftAppUrl, &GiftAppController{})
}

// GiftList App端查询礼物列表(仅已上架,带缓存)
func (c *GiftAppController) GiftList(ctx context.Context, req *giftdto.AppGiftListReq) (res *giftdto.AppGiftListRes, err error) {
	return gift.GetAppGiftList(ctx, req)
}

// ReceivedGiftIds App端查询主播收到过的礼物ID与数量(读缓存,去重最多2000种)
func (c *GiftAppController) ReceivedGiftIds(ctx context.Context, req *giftdto.ReceivedGiftIdsReq) (res *giftdto.ReceivedGiftIdsRes, err error) {
	return gift.GetReceivedGiftIds(ctx, req)
}
