package ipgeo

import "strings"

// Cloudflare CF-IPCountry 特殊值:未知/Tor,不当作有效国家码.
var cfIPCountrySkip = map[string]struct{}{
	"XX": {},
	"T1": {},
}

// CountryNameFromCode 将 ISO 3166-1 alpha-2 转为中文国名;无效或未知返回空.
func CountryNameFromCode(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return ""
	}
	if _, skip := cfIPCountrySkip[code]; skip {
		return ""
	}
	if name, ok := countryNamesZH[code]; ok {
		return name
	}
	return ""
}

// ResolveCountryName 优先 CF-IPCountry(中文),否则按 IP 查 GeoLite(中文优先).
func ResolveCountryName(ip, cfIPCountry string) string {
	if name := CountryNameFromCode(cfIPCountry); name != "" {
		return name
	}
	return GetCountryName(ip)
}
