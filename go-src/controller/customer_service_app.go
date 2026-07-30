package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/customerservicedto"
	"xr-game-server/module/customerservice"
)

const CustomerServiceAppUrl = "/customerService"

type CustomerServiceAppController struct{}

func initCustomerServiceAppController() {
	httpserver.RegNonAuthAPI(CustomerServiceAppUrl, &CustomerServiceAppController{})
}

func (c *CustomerServiceAppController) Cfg(ctx context.Context, req *customerservicedto.AppCustomerServiceCfgReq) (*customerservicedto.AppCustomerServiceCfgRes, error) {
	return customerservice.GetAppCustomerServiceCfg(ctx, req)
}
