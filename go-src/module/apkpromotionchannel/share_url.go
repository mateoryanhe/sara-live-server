package apkpromotionchannel

import (
	"fmt"
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dto/officialsitedeploydto"
)

const (
	// ApkChannelShortPathPattern 根路径短链: GET /{channelCode}
	ApkChannelShortPathPattern = "/{channelCode}"
	// ApkChannelLegacyShortPathPattern 旧版短链(兼容): /dl/apk/{channelCode}
	ApkChannelLegacyShortPathPattern = "/dl/apk/{channelCode}"
)

// ShareShortLinkBaseHost 短链域名(通常与官网同域,网关转发到 Go).
func ShareShortLinkBaseHost() string {
	return strings.TrimSpace(cfg.GetStaticSiteDomain(officialsitedeploydto.OfficialSiteStaticPrefix))
}

// OfficialSiteShareBaseHost CMS 短链预览域名.
func OfficialSiteShareBaseHost() string {
	return ShareShortLinkBaseHost()
}

// BuildChannelShortURL 渠道最短短链: https://{官网域名}/{渠道码}
func BuildChannelShortURL(channelCode string) string {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return ""
	}
	domain := ShareShortLinkBaseHost()
	if domain == "" {
		return ""
	}
	return (&url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   "/" + url.PathEscape(code),
	}).String()
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

func FormatShortLinkPath(channelCode string) string {
	code, ok := normalizeChannelCode(channelCode)
	if !ok {
		return ""
	}
	return fmt.Sprintf("/%s", url.PathEscape(code))
}
