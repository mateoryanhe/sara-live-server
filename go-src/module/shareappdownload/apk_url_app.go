package shareappdownload

import (
	"context"

	"xr-game-server/dto/shareappdownloaddto"
)

// GetApkDownloadUrl App 查询当前 APK 下载链接(内存中官网最大 apk).
func GetApkDownloadUrl(_ context.Context, _ *shareappdownloaddto.GetApkDownloadUrlReq) (*shareappdownloaddto.GetApkDownloadUrlRes, error) {
	return &shareappdownloaddto.GetApkDownloadUrlRes{
		DownloadUrl: GetDownloadUrl(),
	}, nil
}
