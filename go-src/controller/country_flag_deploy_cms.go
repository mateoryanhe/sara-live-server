package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/countryflagdeploydto"
	"xr-game-server/module/countryflagdeploy"
)

const CountryFlagDeployCMSUrl = "/countryFlagDeploy"

type CountryFlagDeployCMSController struct{}

func initCountryFlagDeployCMSController() {
	httpserver.RegCMSHandler(CountryFlagDeployCMSUrl, "/deployZip", handleCountryFlagDeployZip)
	httpserver.RegCMS(CountryFlagDeployCMSUrl, &CountryFlagDeployCMSController{})
}

func handleCountryFlagDeployZip(r *ghttp.Request) {
	res, err := countryflagdeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *CountryFlagDeployCMSController) GetCountryFlagDeployInfo(ctx context.Context, req *countryflagdeploydto.GetCountryFlagDeployInfoReq) (*countryflagdeploydto.GetCountryFlagDeployInfoRes, error) {
	return countryflagdeploy.GetCountryFlagDeployInfo(ctx, req)
}
