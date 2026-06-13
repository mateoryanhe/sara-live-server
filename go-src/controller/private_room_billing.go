package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/privateroombillingdto"
	"xr-game-server/module/privateroombilling"
)

const PrivateRoomBillingUrl = "/privateRoomBilling"

type PrivateRoomBillingController struct{}

func initPrivateRoomBillingController() {
	httpserver.RegCMS(PrivateRoomBillingUrl, &PrivateRoomBillingController{})
}

func (c *PrivateRoomBillingController) BillingList(ctx context.Context, req *privateroombillingdto.BillingListReq) (*httpserver.CMSQueryResp, error) {
	return privateroombilling.GetBillingList(ctx, req)
}

func (c *PrivateRoomBillingController) CreateBilling(ctx context.Context, req *privateroombillingdto.CreateBillingReq) (*privateroombillingdto.CreateBillingRes, error) {
	return privateroombilling.CreateBilling(ctx, req)
}

func (c *PrivateRoomBillingController) UpdateBilling(ctx context.Context, req *privateroombillingdto.UpdateBillingReq) (*privateroombillingdto.UpdateBillingRes, error) {
	return privateroombilling.UpdateBilling(ctx, req)
}

func (c *PrivateRoomBillingController) DeleteBilling(ctx context.Context, req *privateroombillingdto.DeleteBillingReq) (*privateroombillingdto.DeleteBillingRes, error) {
	return privateroombilling.DeleteBilling(ctx, req)
}

func (c *PrivateRoomBillingController) OnShelfBilling(ctx context.Context, req *privateroombillingdto.OnShelfBillingReq) (*privateroombillingdto.OnShelfBillingRes, error) {
	return privateroombilling.OnShelfBilling(ctx, req)
}

func (c *PrivateRoomBillingController) OffShelfBilling(ctx context.Context, req *privateroombillingdto.OffShelfBillingReq) (*privateroombillingdto.OffShelfBillingRes, error) {
	return privateroombilling.OffShelfBilling(ctx, req)
}
