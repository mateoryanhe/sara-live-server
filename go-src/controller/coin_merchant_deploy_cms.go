package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/coinmerchantdeploydto"
	"xr-game-server/module/coinmerchantdeploy"
)

const CoinMerchantDeployCMSUrl = "/coinMerchantDeploy"

type CoinMerchantDeployCMSController struct{}

func initCoinMerchantDeployCMSController() {
	httpserver.RegCMSHandler(CoinMerchantDeployCMSUrl, "/deployZip", handleCoinMerchantDeployZip)
	httpserver.RegCMS(CoinMerchantDeployCMSUrl, &CoinMerchantDeployCMSController{})
}

func handleCoinMerchantDeployZip(r *ghttp.Request) {
	res, err := coinmerchantdeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *CoinMerchantDeployCMSController) GetCoinMerchantDeployInfo(ctx context.Context, req *coinmerchantdeploydto.GetCoinMerchantDeployInfoReq) (*coinmerchantdeploydto.GetCoinMerchantDeployInfoRes, error) {
	return coinmerchantdeploy.GetCoinMerchantDeployInfo(ctx, req)
}

func (c *CoinMerchantDeployCMSController) SaveCoinMerchantDeployCfg(ctx context.Context, req *coinmerchantdeploydto.SaveCoinMerchantDeployCfgReq) (*coinmerchantdeploydto.SaveCoinMerchantDeployCfgRes, error) {
	return coinmerchantdeploy.SaveCoinMerchantDeployCfg(ctx, req)
}
