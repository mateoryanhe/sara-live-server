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
