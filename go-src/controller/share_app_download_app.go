package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
)

const ShareAppDownloadAppUrl = "/shareAppDownload"

type ShareAppDownloadAppController struct{}

func initShareAppDownloadAppController() {
	// APK 下载链接免登录(分享落地/未登录也可拉)
	httpserver.RegNonAuthAPI(ShareAppDownloadAppUrl, &ShareAppDownloadAppController{})
}

// GetApkDownloadUrl 查询当前 APK 下载链接
func (c *ShareAppDownloadAppController) GetApkDownloadUrl(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	return shareappdownload.GetApkDownloadUrl(ctx, req)
}
