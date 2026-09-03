package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/thirdpaydeploydto"
	"xr-game-server/module/thirdpaydeploy"
)

const ThirdPayDeployCMSUrl = "/thirdPayDeploy"

type ThirdPayDeployCMSController struct{}

func initThirdPayDeployCMSController() {
	httpserver.RegCMSHandler(ThirdPayDeployCMSUrl, "/deployZip", handleThirdPayDeployZip)
	httpserver.RegCMS(ThirdPayDeployCMSUrl, &ThirdPayDeployCMSController{})
}

func handleThirdPayDeployZip(r *ghttp.Request) {
	res, err := thirdpaydeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *ThirdPayDeployCMSController) GetThirdPayDeployInfo(ctx context.Context, req *thirdpaydeploydto.GetThirdPayDeployInfoReq) (*thirdpaydeploydto.GetThirdPayDeployInfoRes, error) {
	return thirdpaydeploy.GetThirdPayDeployInfo(ctx, req)
}
