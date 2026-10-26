package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/trackingdto"
	"xr-game-server/module/tracking"
)

const TrackingDownloadAppUrl = "/trackingEvent"

type TrackingDownloadAppController struct{}

func initTrackingDownloadAppController() {
	// H5 / 官网下载上报均免登录(下载落地页无 token)
	httpserver.RegNonAuthAPI(TrackingDownloadAppUrl, &TrackingDownloadAppController{})
}

// ReportH5Download 上报 H5 下载点击(免登录)
func (c *TrackingDownloadAppController) ReportH5Download(ctx context.Context, req *trackingdto.ReportH5DownloadReq) (*trackingdto.ReportH5DownloadRes, error) {
	return tracking.ReportH5Download(ctx, req)
}

// ReportOfficialSiteDownload 上报官网下载点击(免登录)
func (c *TrackingDownloadAppController) ReportOfficialSiteDownload(ctx context.Context, req *trackingdto.ReportOfficialSiteDownloadReq) (*trackingdto.ReportOfficialSiteDownloadRes, error) {
	return tracking.ReportOfficialSiteDownload(ctx, req)
}
