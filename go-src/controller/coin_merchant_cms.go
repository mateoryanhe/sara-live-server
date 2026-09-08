package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/coinmerchantdto"
	"xr-game-server/module/coinmerchant"
)

const CoinMerchantCMSUrl = "/coinMerchant"

type CoinMerchantController struct{}

func initCoinMerchantController() {
	httpserver.RegCMS(CoinMerchantCMSUrl, &CoinMerchantController{})
}

func (c *CoinMerchantController) CoinMerchantList(ctx context.Context, req *coinmerchantdto.CoinMerchantListReq) (*httpserver.CMSQueryResp, error) {
	return coinmerchant.ListCoinMerchants(ctx, req)
}

func (c *CoinMerchantController) CreateCoinMerchant(ctx context.Context, req *coinmerchantdto.CreateCoinMerchantReq) (*coinmerchantdto.CreateCoinMerchantRes, error) {
	return coinmerchant.CreateCoinMerchant(ctx, req)
}

func (c *CoinMerchantController) ResetCoinMerchantPassword(ctx context.Context, req *coinmerchantdto.ResetCoinMerchantPasswordReq) (*coinmerchantdto.ResetCoinMerchantPasswordRes, error) {
	return coinmerchant.ResetCoinMerchantPassword(ctx, req)
}

func (c *CoinMerchantController) CancelCoinMerchant(ctx context.Context, req *coinmerchantdto.CancelCoinMerchantReq) (*coinmerchantdto.CancelCoinMerchantRes, error) {
	return coinmerchant.CancelCoinMerchant(ctx, req)
}
