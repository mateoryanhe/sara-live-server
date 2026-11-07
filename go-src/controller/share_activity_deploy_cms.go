package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareactivitydeploydto"
	"xr-game-server/module/shareactivitydeploy"
)

const ShareActivityDeployCMSUrl = "/shareActivityDeploy"

type ShareActivityDeployCMSController struct{}

func initShareActivityDeployCMSController() {
	httpserver.RegCMSHandler(ShareActivityDeployCMSUrl, "/deployZip", handleShareActivityDeployZip)
	httpserver.RegCMS(ShareActivityDeployCMSUrl, &ShareActivityDeployCMSController{})
}

func handleShareActivityDeployZip(r *ghttp.Request) {
	res, err := shareactivitydeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *ShareActivityDeployCMSController) GetShareActivityDeployInfo(ctx context.Context, req *shareactivitydeploydto.GetShareActivityDeployInfoReq) (*shareactivitydeploydto.GetShareActivityDeployInfoRes, error) {
	return shareactivitydeploy.GetShareActivityDeployInfo(ctx, req)
}

func (c *ShareActivityDeployCMSController) SaveShareActivityDeployCfg(ctx context.Context, req *shareactivitydeploydto.SaveShareActivityDeployCfgReq) (*shareactivitydeploydto.SaveShareActivityDeployCfgRes, error) {
	return shareactivitydeploy.SaveShareActivityDeployCfg(ctx, req)
}
