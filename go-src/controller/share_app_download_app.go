package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/shareappdownload"
	"xr-game-server/module/stat"
	"xr-game-server/module/tracking"
)

const ShareAppDownloadAppUrl = "/shareAppDownload"

type ShareAppDownloadAppController struct{}

func initShareAppDownloadAppController() {
	// APK 下载链接免登录(分享落地/未登录也可拉)
	httpserver.RegNonAuthAPI(ShareAppDownloadAppUrl, &ShareAppDownloadAppController{})
	// 短链: 点击后 302 到官网 assets 当前最大 APK
	httpserver.RegNonAuthGETHandler(shareappdownload.ApkShortPath, handleApkShortRedirect)
	// Sofie: Google Play 短链(默认 /dl/sofie; 按包名 /dl/p/{packageName})
	httpserver.RegNonAuthGETHandler(shareappdownload.SofieDefaultShortPath, handleSofieDefaultShortRedirect)
	httpserver.RegNonAuthGETHandler(shareappdownload.SofiePackageShortPathPattern, handleSofiePackageShortRedirect)
}

// GetApkDownloadUrl 查询当前 APK 下载链接(可选 ch 匹配推广渠道包)
func (c *ShareAppDownloadAppController) GetApkDownloadUrl(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	res, err := shareappdownload.GetApkDownloadUrl(ctx, req)
	if err != nil {
		return nil, err
	}
	if code := shareappdownload.LookupApkDownloadStatChannelCode(ctx, req); code > 0 {
		stat.EnqueuePromotionChannelH5DownloadApk(code, time.Now())
	}
	return res, nil
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

func handleSofieDefaultShortRedirect(r *ghttp.Request) {
	target := shareappdownload.ResolveSofieGooglePlayUrl("")
	if target == "" {
		r.Response.WriteStatus(http.StatusNotFound, "not available")
		return
	}
	tracking.TrackShareShortLinkClick()
	r.Response.RedirectTo(target, http.StatusFound)
}

func handleSofiePackageShortRedirect(r *ghttp.Request) {
	pkg := shareappdownload.ParseSofiePackageFromShortPath(r.Get("packageName").String())
	target := shareappdownload.ResolveSofieGooglePlayUrl(pkg)
	if target == "" {
		r.Response.WriteStatus(http.StatusNotFound, "not available")
		return
	}
	tracking.TrackShareShortLinkClick()
	r.Response.RedirectTo(target, http.StatusFound)
}
