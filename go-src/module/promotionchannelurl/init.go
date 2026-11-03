package promotionchannelurl

import (
	"strings"

	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/module/upload"
)

// Init 预热推广渠道 APK URL 缓存.
func Init() {
	WarmupURLCacheFromDB()
}

func WarmupURLCacheFromDB() {
	for _, row := range promotionchannelurldao.ListAllForURLCache() {
		if row == nil || row.ChannelCode <= 0 {
			continue
		}
		url := strings.TrimSpace(row.ApkDownloadURL)
		if url == "" {
			url = upload.PublicURLForPromotionChannelUrlApkStoredName(row.ApkFileName)
		}
		SetCachedApkURL(row.ChannelCode, url)
	}
}
