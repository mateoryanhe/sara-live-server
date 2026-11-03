package promotionchannelurl

import (
	"strconv"
	"strings"
)

// ParseChannelCodeFromQuery lookupApkDownload / 落地页 ch: 推广渠道 URL 为整数编码.
func ParseChannelCodeFromQuery(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	code, err := strconv.Atoi(raw)
	if err != nil || code <= 0 {
		return 0, false
	}
	return code, true
}

// ResolveApkDownloadForChannelQuery 带 ch 查 APK: 命中推广渠道 URL 包则返回渠道 URL,否则 officialFallback.
// 第二返回值 >0 表示命中渠道(可用于 H5 下载 APK 统计).
func ResolveApkDownloadForChannelQuery(rawCh, officialFallback string) (downloadURL string, statChannelCode int) {
	downloadURL = officialFallback
	code, ok := ParseChannelCodeFromQuery(rawCh)
	if !ok || !HasCachedChannel(code) {
		return downloadURL, 0
	}
	channelURL := GetCachedApkURL(code)
	if channelURL == "" {
		return downloadURL, 0
	}
	return channelURL, code
}
