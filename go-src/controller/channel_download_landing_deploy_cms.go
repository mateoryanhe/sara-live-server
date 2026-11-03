package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/channeldownloadlandingdeploydto"
	"xr-game-server/module/channeldownloadlandingdeploy"
)

const ChannelDownloadLandingDeployCMSUrl = "/channelDownloadLandingDeploy"

type ChannelDownloadLandingDeployCMSController struct{}

func initChannelDownloadLandingDeployCMSController() {
	httpserver.RegCMSHandler(ChannelDownloadLandingDeployCMSUrl, "/deployZip", handleChannelDownloadLandingDeployZip)
	httpserver.RegCMS(ChannelDownloadLandingDeployCMSUrl, &ChannelDownloadLandingDeployCMSController{})
}

func handleChannelDownloadLandingDeployZip(r *ghttp.Request) {
	res, err := channeldownloadlandingdeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *ChannelDownloadLandingDeployCMSController) GetChannelDownloadLandingDeployInfo(ctx context.Context, req *channeldownloadlandingdeploydto.GetChannelDownloadLandingDeployInfoReq) (*channeldownloadlandingdeploydto.GetChannelDownloadLandingDeployInfoRes, error) {
	return channeldownloadlandingdeploy.GetChannelDownloadLandingDeployInfo(ctx, req)
}

func (c *ChannelDownloadLandingDeployCMSController) SaveChannelDownloadLandingDeployCfg(ctx context.Context, req *channeldownloadlandingdeploydto.SaveChannelDownloadLandingDeployCfgReq) (*channeldownloadlandingdeploydto.SaveChannelDownloadLandingDeployCfgRes, error) {
	return channeldownloadlandingdeploy.SaveChannelDownloadLandingDeployCfg(ctx, req)
}
