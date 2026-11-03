package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/apkpromotionchanneldto"
	"xr-game-server/module/apkpromotionchannel"
	"xr-game-server/module/officialsitedeploy"
	"xr-game-server/module/shareappdownload"
	"xr-game-server/module/stat"
)

const ApkPromotionChannelCMSUrl = "/apkPromotionChannel"

type ApkPromotionChannelCMSController struct{}

func initApkPromotionChannelCMSController() {
	httpserver.RegCMS(ApkPromotionChannelCMSUrl, &ApkPromotionChannelCMSController{})
	httpserver.RegCMSHandler(ApkPromotionChannelCMSUrl, "/uploadApkFileChunk", handleApkPromotionChannelUploadApkChunk)
	initApkPromotionChannelShortLinkHandler()
}

func handleApkPromotionChannelUploadApkChunk(r *ghttp.Request) {
	res, err := officialsitedeploy.UploadApkPromotionChannelApkChunkFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *ApkPromotionChannelCMSController) ApkPromotionChannelList(ctx context.Context, req *apkpromotionchanneldto.ApkPromotionChannelListReq) (*apkpromotionchanneldto.ApkPromotionChannelListRes, error) {
	return apkpromotionchannel.GetList(ctx, req)
}

func (c *ApkPromotionChannelCMSController) CreateApkPromotionChannel(ctx context.Context, req *apkpromotionchanneldto.CreateApkPromotionChannelReq) (*apkpromotionchanneldto.CreateApkPromotionChannelRes, error) {
	return apkpromotionchannel.Create(ctx, req)
}

func (c *ApkPromotionChannelCMSController) UpdateApkPromotionChannel(ctx context.Context, req *apkpromotionchanneldto.UpdateApkPromotionChannelReq) (*apkpromotionchanneldto.UpdateApkPromotionChannelRes, error) {
	return apkpromotionchannel.Update(ctx, req)
}

func (c *ApkPromotionChannelCMSController) DeleteApkPromotionChannel(ctx context.Context, req *apkpromotionchanneldto.DeleteApkPromotionChannelReq) (*apkpromotionchanneldto.DeleteApkPromotionChannelRes, error) {
	return apkpromotionchannel.Delete(ctx, req)
}

func (c *ApkPromotionChannelCMSController) InitApkPromotionChannelApkUpload(ctx context.Context, req *apkpromotionchanneldto.InitApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.InitApkPromotionChannelApkUploadRes, error) {
	return officialsitedeploy.InitApkPromotionChannelApkUpload(ctx, req)
}

func (c *ApkPromotionChannelCMSController) CompleteApkPromotionChannelApkUpload(ctx context.Context, req *apkpromotionchanneldto.CompleteApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.CompleteApkPromotionChannelApkUploadRes, error) {
	return officialsitedeploy.CompleteApkPromotionChannelApkUpload(ctx, req)
}

func (c *ApkPromotionChannelCMSController) AbortApkPromotionChannelApkUpload(ctx context.Context, req *apkpromotionchanneldto.AbortApkPromotionChannelApkUploadReq) (*apkpromotionchanneldto.AbortApkPromotionChannelApkUploadRes, error) {
	return officialsitedeploy.AbortApkPromotionChannelApkUpload(ctx, req)
}

func (c *ApkPromotionChannelCMSController) ApkPromotionChannelDailyStatList(ctx context.Context, req *apkpromotionchanneldto.ApkPromotionChannelDailyStatListReq) (*apkpromotionchanneldto.ApkPromotionChannelDailyStatListRes, error) {
	return apkpromotionchannel.GetDailyStatList(ctx, req)
}

func initApkPromotionChannelShortLinkHandler() {
	// 渠道落地页: GET ?ch= 上报并 302(免登录,不走 H5 加解密)
	httpserver.RegNonAuthGETHandler(apkpromotionchannel.ApkChannelPublicDownloadPath, handleApkPromotionChannelDownloadByQuery)
	// 先注册旧路径兼容;根路径 /{channelCode} 仅对已注册渠道 302(见 handler).
	httpserver.RegNonAuthGETHandler(apkpromotionchannel.ApkChannelLegacyShortPathPattern, handleApkPromotionChannelShortRedirect)
	httpserver.RegNonAuthGETHandler(apkpromotionchannel.ApkChannelShortPathPattern, handleApkPromotionChannelShortRedirect)
}

func handleApkPromotionChannelDownloadByQuery(r *ghttp.Request) {
	code, ok := apkpromotionchannel.ParseChannelCodeFromQuery(
		r.GetQuery(apkpromotionchannel.ChannelDownloadQueryParam).String(),
		r.GetQuery("channelCode").String(),
		r.GetQuery("channel_code").String(),
	)
	if !ok {
		redirectOfficialSiteApkOrNotFound(r)
		return
	}
	apkPromotionChannelDownloadRedirectWithOfficialFallback(r, code)
}

func handleApkPromotionChannelShortRedirect(r *ghttp.Request) {
	code, ok := apkpromotionchannel.ParseChannelCodeFromRequest(r.Get("channelCode").String())
	if !ok {
		r.Response.WriteStatus(http.StatusNotFound, "not found")
		return
	}
	apkPromotionChannelDownloadRedirect(r, code)
}

// apkPromotionChannelDownloadRedirectWithOfficialFallback 落地页 ?ch= 上报: 渠道无效或未配置 APK 时回落官网包.
func apkPromotionChannelDownloadRedirectWithOfficialFallback(r *ghttp.Request, code string) {
	if !apkpromotionchannel.HasCachedChannel(code) {
		redirectOfficialSiteApkOrNotFound(r)
		return
	}
	target := apkpromotionchannel.ResolveChannelApkDownloadURL(code)
	if target == "" {
		redirectOfficialSiteApkOrNotFound(r)
		return
	}
	stat.EnqueueApkPromotionChannelDownloadClick(code, time.Now())
	r.Response.RedirectTo(target, http.StatusFound)
}

func redirectOfficialSiteApkOrNotFound(r *ghttp.Request) {
	target := shareappdownload.GetDownloadUrl()
	if target == "" {
		r.Response.WriteStatus(http.StatusNotFound, "apk not available")
		return
	}
	r.Response.RedirectTo(target, http.StatusFound)
}

func apkPromotionChannelDownloadRedirect(r *ghttp.Request, code string) {
	if !apkpromotionchannel.HasCachedChannel(code) {
		r.Response.WriteStatus(http.StatusNotFound, "not found")
		return
	}
	target := apkpromotionchannel.ResolveChannelApkDownloadURL(code)
	if target == "" {
		r.Response.WriteStatus(http.StatusNotFound, "apk not available")
		return
	}
	stat.EnqueueApkPromotionChannelDownloadClick(code, time.Now())
	r.Response.RedirectTo(target, http.StatusFound)
}
