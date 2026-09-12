package ipgeo

import (
	"strings"

	"xr-game-server/constants/country"
)

// Cloudflare CF-IPCountry 特殊值:未知/Tor,不当作有效国家码.
var cfIPCountrySkip = map[string]struct{}{
	"XX": {},
	"T1": {},
}

// NormalizeCountryCode 规范化 ISO 3166-1 alpha-2;无效/特殊值返回空.
func NormalizeCountryCode(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 2 {
		return ""
	}
	if _, skip := cfIPCountrySkip[code]; skip {
		return ""
	}
	return code
}

// ResolveCountryCode 优先 CF-IPCountry 简码,否则 GeoLite ISO 码;入库统一存 Code.
func ResolveCountryCode(ip, cfIPCountry string) string {
	if code := NormalizeCountryCode(cfIPCountry); code != "" {
		return code
	}
	return NormalizeCountryCode(GetCountryCode(ip))
}

// FormatCountryDisplay 将入库简码格式化为「中文 / English」(CMS 展示).
func FormatCountryDisplay(codeOrLegacy string) string {
	return country.FormatZhEn(codeOrLegacy)
}
