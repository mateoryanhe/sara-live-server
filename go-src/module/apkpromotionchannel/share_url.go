package apkpromotionchannel

import (
	"fmt"
	"net/url"
	"strings"

	"xr-game-server/dto/channeldownloadlandingdeploydto"
	"xr-game-server/module/channeldownloadlandingdeploy"
)

const (
	// ApkChannelShortPathPattern 根路径短链: GET /{channelCode}
	ApkChannelShortPathPattern = "/{channelCode}"
	// ApkChannelLegacyShortPathPattern 旧版短链(兼容): /dl/apk/{channelCode}
	ApkChannelLegacyShortPathPattern = "/dl/apk/{channelCode}"
	// ApkChannelPublicDownloadPath 渠道落地页免登录上报: GET ?ch=渠道码 → 302 APK
	ApkChannelPublicDownloadPath = "/apkPromotionChannel/download"
)

// ChannelDownloadQueryParam 落地页渠道码 query 参数(优先 ch).
const ChannelDownloadQueryParam = "ch"

// ShareShortLinkBaseHost 推广短链域名(渠道下载落地页站点).
func ShareShortLinkBaseHost() string {
	return channeldownloadlandingdeploy.ShareBaseHost()
}

// OfficialSiteShareBaseHost CMS 列表 shareBaseHost 字段(历史命名,实为落地页域名).
func OfficialSiteShareBaseHost() string {
	return ShareShortLinkBaseHost()
}

// BuildChannelShortURL 推广短链: https://{落地页域名}/channel-download/?ch={渠道码}
func BuildChannelShortURL(channelCode string) string {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return ""
	}
	return channeldownloadlandingdeploy.BuildShareURL(code)
}

// ResolveChannelApkDownloadURL 短链 302 目标:只读内存缓存,不查库.
func ResolveChannelApkDownloadURL(channelCode string) string {
	return GetCachedApkURL(channelCode)
}

// BuildChannelApkDirectDownloadURL CMS 展示用(与内存缓存一致).
func BuildChannelApkDirectDownloadURL(channelCode string) string {
	return ResolveChannelApkDownloadURL(channelCode)
}

func ParseChannelCodeFromRequest(raw string) (string, bool) {
	decoded, err := url.PathUnescape(strings.TrimSpace(raw))
	if err != nil {
		decoded = strings.TrimSpace(raw)
	}
	return normalizeChannelCode(decoded)
}

// ParseChannelCodeFromQuery 落地页 query: ch / channelCode / channel_code.
func ParseChannelCodeFromQuery(ch, channelCode, channelCodeSnake string) (string, bool) {
	for _, raw := range []string{ch, channelCode, channelCodeSnake} {
		if code, ok := ParseChannelCodeFromRequest(raw); ok {
			return code, true
		}
	}
	return "", false
}

func FormatShortLinkPath(channelCode string) string {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return ""
	}
	prefix := strings.TrimSuffix(channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix, "/")
	return fmt.Sprintf("%s/?%s=%s", prefix, ChannelDownloadQueryParam, url.QueryEscape(code))
}
