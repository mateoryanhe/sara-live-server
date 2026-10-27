package controller

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
	"xr-game-server/module/tracking"
)

const ShareAppDownloadAppUrl = "/shareAppDownload"

type ShareAppDownloadAppController struct{}

func initShareAppDownloadAppController() {
	// APK 下载链接免登录(分享落地/未登录也可拉)
	httpserver.RegNonAuthAPI(ShareAppDownloadAppUrl, &ShareAppDownloadAppController{})
	// 短链: 点击后 302 到官网 assets 当前最大 APK
	httpserver.RegNonAuthGETHandler(shareappdownload.ApkShortPath, handleApkShortRedirect)
}

// GetApkDownloadUrl 查询当前 APK 下载链接
func (c *ShareAppDownloadAppController) GetApkDownloadUrl(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	return shareappdownload.GetApkDownloadUrl(ctx, req)
}

func handleApkShortRedirect(r *ghttp.Request) {
	target := shareappdownload.GetDownloadUrl()
	if target == "" {
		r.Response.WriteStatus(http.StatusNotFound, "apk not available")
		return
	}
	// 成功 302 前计数(分享短链点击量)
	tracking.TrackShareShortLinkClick()
	r.Response.RedirectTo(target, http.StatusFound)
}
