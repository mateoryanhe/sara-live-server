package promotionchannelurl

import (
	"fmt"
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/h5livedeploydto"
)

const promotionShareQueryKey = "ch"

// H5LiveShareBaseHost 从 domain_site_mappings 读取 H5 直播站点域名(直查库,与 CMS H5 部署页一致)。
func H5LiveShareBaseHost() string {
	row := domainsitedao.GetMappingBySiteKey(h5livedeploydto.H5LiveSiteKey)
	if row == nil {
		return ""
	}
	domains := cfg.SplitDomains(strings.TrimSpace(row.Domain))
	if len(domains) == 0 {
		return ""
	}
	return domains[0]
}

// BuildShareURL 生成推广分享链接,如 https://h5-live.example.com?ch=1
func BuildShareURL(channelCode int) string {
	if channelCode <= 0 {
		return ""
	}
	host := H5LiveShareBaseHost()
	if host == "" {
		return ""
	}
	u := url.URL{
		Scheme:   "https",
		Host:     host,
		RawQuery: fmt.Sprintf("%s=%d", promotionShareQueryKey, channelCode),
	}
	return u.String()
}
