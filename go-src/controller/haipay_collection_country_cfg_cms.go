package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/haipaydto"
	"xr-game-server/module/recharge"
)

const HaiPayCollectionCountryCfgCMSURL = "/paymentCountryCfg"

type HaiPayCollectionCountryCfgCMSController struct{}

func initPaymentCountryCfgCMSController() {
	httpserver.RegCMS(HaiPayCollectionCountryCfgCMSURL, &HaiPayCollectionCountryCfgCMSController{})
}

func (c *HaiPayCollectionCountryCfgCMSController) GetCollectionCountryCfg(ctx context.Context, req *haipaydto.GetCollectionCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	return recharge.GetHaiPayCollectionCountryCfg(ctx, req)
}

func (c *HaiPayCollectionCountryCfgCMSController) SaveCollectionCountryCfg(ctx context.Context, req *haipaydto.SaveCollectionCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	return recharge.SaveHaiPayCollectionCountryCfg(ctx, req)
}
