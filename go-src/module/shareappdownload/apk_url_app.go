package shareappdownload

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/promotionchannelurl"
)

// GetApkDownloadUrl 查询 APK 下载链接: 可选 ch(推广渠道 URL 整数编码) 匹配渠道包,否则回落官网最大 apk.
// Sofie 路径按 X-App-Package 匹配 Google Play 配置,未命中则用默认行.
func GetApkDownloadUrl(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	downloadURL, _ := lookupApkDownload(ctx, req)
	return &shareappdownloaddto.GetApkDownloadUrlRes{
		DownloadUrl: downloadURL,
	}, nil
}

// LookupApkDownloadStatChannelCode 与 GetApkDownloadUrl 同规则; >0 时计入推广渠道 H5 下载 APK 统计.
func LookupApkDownloadStatChannelCode(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) int {
	_, statCode := lookupApkDownload(ctx, req)
	return statCode
}

func lookupApkDownload(ctx context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (downloadURL string, statChannelCode int) {
	if httpserver.IsSofieAppRequest(ctx) {
		pkg := httpserver.GetPackageNameFromContext(ctx)
		if url := ResolveSofieGooglePlayUrl(pkg); url != "" {
			return url, 0
		}
	}
	ch := ""
	if req != nil {
		ch = req.Ch
	}
	return promotionchannelurl.ResolveApkDownloadForChannelQuery(ch, GetDownloadUrl())
}
