package shareappdownload

import (
	"context"

	"xr-game-server/dto/shareappdownloaddto"
	"xr-game-server/module/promotionchannelurl"
)

// GetApkDownloadUrl 查询 APK 下载链接: 可选 ch(推广渠道 URL 整数编码) 匹配渠道包,否则回落官网最大 apk.
func GetApkDownloadUrl(_ context.Context, req *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	downloadURL, _ := lookupApkDownload(req)
	return &shareappdownloaddto.GetApkDownloadUrlRes{
		DownloadUrl: downloadURL,
	}, nil
}

// LookupApkDownloadStatChannelCode 与 GetApkDownloadUrl 同规则; >0 时计入推广渠道 H5 下载 APK 统计.
func LookupApkDownloadStatChannelCode(req *shareappdownloaddto.GetApkDownloadUrlReq) int {
	_, statCode := lookupApkDownload(req)
	return statCode
}

func lookupApkDownload(req *shareappdownloaddto.GetApkDownloadUrlReq) (downloadURL string, statChannelCode int) {
	ch := ""
	if req != nil {
		ch = req.Ch
	}
	return promotionchannelurl.ResolveApkDownloadForChannelQuery(ch, GetDownloadUrl())
}
