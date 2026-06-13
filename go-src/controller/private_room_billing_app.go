package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/privateroombillingdto"
	"xr-game-server/module/privateroombilling"
)

const PrivateRoomBillingAppUrl = "/privateRoomBilling"

type PrivateRoomBillingAppController struct{}

func initPrivateRoomBillingAppController() {
	httpserver.RegAPI(PrivateRoomBillingAppUrl, &PrivateRoomBillingAppController{})
}

func (c *PrivateRoomBillingAppController) AppBillingList(ctx context.Context, req *privateroombillingdto.AppBillingListReq) (*privateroombillingdto.AppBillingListRes, error) {
	return privateroombilling.GetAppBillingList(ctx, req)
}
