package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/h5livedeploydto"
	"xr-game-server/module/h5livedeploy"
)

const H5LiveDeployCMSUrl = "/h5LiveDeploy"

type H5LiveDeployCMSController struct{}

func initH5LiveDeployCMSController() {
	httpserver.RegCMSHandler(H5LiveDeployCMSUrl, "/deployZip", handleH5LiveDeployZip)
	httpserver.RegCMS(H5LiveDeployCMSUrl, &H5LiveDeployCMSController{})
}

func handleH5LiveDeployZip(r *ghttp.Request) {
	res, err := h5livedeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *H5LiveDeployCMSController) GetH5LiveDeployInfo(ctx context.Context, req *h5livedeploydto.GetH5LiveDeployInfoReq) (*h5livedeploydto.GetH5LiveDeployInfoRes, error) {
	return h5livedeploy.GetH5LiveDeployInfo(ctx, req)
}

func (c *H5LiveDeployCMSController) SaveH5LiveDeployCfg(ctx context.Context, req *h5livedeploydto.SaveH5LiveDeployCfgReq) (*h5livedeploydto.SaveH5LiveDeployCfgRes, error) {
	return h5livedeploy.SaveH5LiveDeployCfg(ctx, req)
}
