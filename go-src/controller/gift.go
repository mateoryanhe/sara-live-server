package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/giftdto"
	"xr-game-server/module/gift"
)

const (
	GiftUrl = "/gift"
)

type GiftController struct{}

func initGiftController() {
	httpserver.RegCMS(GiftUrl, &GiftController{})
}

// GiftList 分页获取礼物列表
func (c *GiftController) GiftList(ctx context.Context, req *giftdto.GiftListReq) (res *httpserver.CMSQueryResp, err error) {
	return gift.GetGiftList(ctx, req)
}

// CreateGift 创建礼物
func (c *GiftController) CreateGift(ctx context.Context, req *giftdto.CreateGiftReq) (res *giftdto.CreateGiftRes, err error) {
	return gift.CreateGift(ctx, req)
}

// UpdateGift 修改礼物
func (c *GiftController) UpdateGift(ctx context.Context, req *giftdto.UpdateGiftReq) (res *giftdto.UpdateGiftRes, err error) {
	return gift.UpdateGift(ctx, req)
}

// DeleteGift 删除礼物
func (c *GiftController) DeleteGift(ctx context.Context, req *giftdto.DeleteGiftReq) (res *giftdto.DeleteGiftRes, err error) {
	return gift.DeleteGift(ctx, req)
}

// OnShelfGift 上架礼物
func (c *GiftController) OnShelfGift(ctx context.Context, req *giftdto.OnShelfGiftReq) (res *giftdto.OnShelfGiftRes, err error) {
	return gift.OnShelfGift(ctx, req)
}

// OffShelfGift 下架礼物
func (c *GiftController) OffShelfGift(ctx context.Context, req *giftdto.OffShelfGiftReq) (res *giftdto.OffShelfGiftRes, err error) {
	return gift.OffShelfGift(ctx, req)
}
