package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
)

const ShareAppDownloadCMSUrl = "/shareAppDownload"

type ShareAppDownloadCMSController struct{}

func initShareAppDownloadCMSController() {
	httpserver.RegCMS(ShareAppDownloadCMSUrl, &ShareAppDownloadCMSController{})
}

func (c *ShareAppDownloadCMSController) GetShareAppDownloadCfg(ctx context.Context, req *shareappdownloaddto.GetShareAppDownloadCfgReq) (*shareappdownloaddto.GetShareAppDownloadCfgRes, error) {
	return shareappdownload.GetShareAppDownloadCfg(ctx, req)
}

func (c *ShareAppDownloadCMSController) SaveShareAppDownloadCfg(ctx context.Context, req *shareappdownloaddto.SaveShareAppDownloadCfgReq) (*shareappdownloaddto.SaveShareAppDownloadCfgRes, error) {
	return shareappdownload.SaveShareAppDownloadCfg(ctx, req)
}
