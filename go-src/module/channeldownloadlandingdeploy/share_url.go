package channeldownloadlandingdeploy

import (
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/channeldownloadlandingdeploydto"
)

// ShareQueryKey 落地页渠道码 query 参数(与 APK 推广 API ?ch= 一致).
const ShareQueryKey = "ch"

// ShareBaseHost 渠道下载落地页域名(与 CMS「渠道下载落地页部署」、domain_site_mappings 一致).
func ShareBaseHost() string {
	row := domainsitedao.GetMappingBySiteKey(channeldownloadlandingdeploydto.ChannelDownloadLandingSiteKey)
	if row != nil {
		domains := cfg.SplitDomains(strings.TrimSpace(row.Domain))
		if len(domains) > 0 {
			return domains[0]
		}
	}
	return strings.TrimSpace(cfg.GetStaticSiteDomain(channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix))
}

// BuildShareURL 推广短链: https://{域名}/channel-download/?ch={渠道码}
func BuildShareURL(channelCode string) string {
	channelCode = strings.TrimSpace(channelCode)
	if channelCode == "" {
		return ""
	}
	host := ShareBaseHost()
	if host == "" {
		return ""
	}
	prefix := strings.TrimSuffix(channeldownloadlandingdeploydto.ChannelDownloadLandingStaticPrefix, "/")
	if prefix == "" {
		prefix = "/"
	}
	u := url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     prefix + "/",
		RawQuery: url.Values{ShareQueryKey: {channelCode}}.Encode(),
	}
	return u.String()
}
