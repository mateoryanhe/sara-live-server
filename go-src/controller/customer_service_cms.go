package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/customerservicedto"
	"xr-game-server/module/customerservice"
)

const CustomerServiceCMSUrl = "/customerService"

type CustomerServiceCMSController struct{}

func initCustomerServiceCMSController() {
	httpserver.RegCMS(CustomerServiceCMSUrl, &CustomerServiceCMSController{})
}

func (c *CustomerServiceCMSController) GetCustomerServiceCfg(ctx context.Context, req *customerservicedto.GetCustomerServiceCfgReq) (*customerservicedto.GetCustomerServiceCfgRes, error) {
	return customerservice.GetCustomerServiceCfg(ctx, req)
}

func (c *CustomerServiceCMSController) SaveCustomerServiceCfg(ctx context.Context, req *customerservicedto.SaveCustomerServiceCfgReq) (*customerservicedto.SaveCustomerServiceCfgRes, error) {
	return customerservice.SaveCustomerServiceCfg(ctx, req)
}
