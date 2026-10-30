package apkpromotionchannel

import (
	"strings"

	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/module/upload"
)

// Init 进程内预热渠道 APK URL 缓存(启动/热更后一次读库).
func Init() {
	WarmupURLCacheFromDB()
}

func WarmupURLCacheFromDB() {
	rows := apkpromotionchanneldao.ListAllForURLCache()
	for _, row := range rows {
		if row == nil {
			continue
		}
		code := strings.TrimSpace(row.ChannelCode)
		if code == "" {
			continue
		}
		url := strings.TrimSpace(row.ApkDownloadURL)
		if url == "" {
			url = upload.PublicURLForApkPromotionStoredName(row.ApkFileName)
		}
		SetCachedApkURL(code, url)
	}
}
