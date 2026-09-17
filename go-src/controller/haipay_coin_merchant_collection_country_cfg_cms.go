package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/haipaydto"
	"xr-game-server/module/recharge"
)

const HaiPayCoinMerchantCollectionCountryCfgCMSURL = "/coinMerchantPaymentCountryCfg"

type HaiPayCoinMerchantCollectionCountryCfgCMSController struct{}

func initCoinMerchantPaymentCountryCfgCMSController() {
	httpserver.RegCMS(HaiPayCoinMerchantCollectionCountryCfgCMSURL, &HaiPayCoinMerchantCollectionCountryCfgCMSController{})
}

func (c *HaiPayCoinMerchantCollectionCountryCfgCMSController) GetCollectionCountryCfg(ctx context.Context, req *haipaydto.GetCollectionCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	return recharge.GetHaiPayCoinMerchantCollectionCountryCfg(ctx, req)
}

func (c *HaiPayCoinMerchantCollectionCountryCfgCMSController) SaveCollectionCountryCfg(ctx context.Context, req *haipaydto.SaveCoinMerchantCollectionCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	return recharge.SaveHaiPayCoinMerchantCollectionCountryCfg(ctx, req)
}
