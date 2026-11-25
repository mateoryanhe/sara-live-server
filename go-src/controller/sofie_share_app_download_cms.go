package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
)

const SofieShareAppDownloadCMSUrl = "/sofieShareAppDownload"

type SofieShareAppDownloadCMSController struct{}

func initSofieShareAppDownloadCMSController() {
	httpserver.RegCMS(SofieShareAppDownloadCMSUrl, &SofieShareAppDownloadCMSController{})
}

func (c *SofieShareAppDownloadCMSController) GetSofieShareAppDownloadCfg(ctx context.Context, req *shareappdownloaddto.GetSofieShareAppDownloadCfgReq) (*shareappdownloaddto.GetSofieShareAppDownloadCfgRes, error) {
	return shareappdownload.GetSofieShareAppDownloadCfg(ctx, req)
}

func (c *SofieShareAppDownloadCMSController) SaveSofieShareAppDownloadCfg(ctx context.Context, req *shareappdownloaddto.SaveSofieShareAppDownloadCfgReq) (*shareappdownloaddto.SaveSofieShareAppDownloadCfgRes, error) {
	return shareappdownload.SaveSofieShareAppDownloadCfg(ctx, req)
}