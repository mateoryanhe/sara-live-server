package attrreport

import (
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/module/ipgeo"
)

// ResolveCountryForUser 归因上报 country.
// 与全站一致: constants/country ISO 3166-1 alpha-2(如 ID=印尼);
// 优先注册国 → 登录国 → IP Geo; 仅输出规范化简码.
func ResolveCountryForUser(userId uint64, clientIP string) string {
	if userId > 0 {
		if loc := userloginlocationdao.GetByUserId(userId); loc != nil {
			if c := unifyCountryCode(loc.RegisterCountry); c != "" {
				return c
			}
			if c := unifyCountryCode(loc.LoginCountry); c != "" {
				return c
			}
			if clientIP == "" {
				clientIP = strings.TrimSpace(loc.RegisterIp)
			}
			if clientIP == "" {
				clientIP = strings.TrimSpace(loc.IP)
			}
		}
	}
	clientIP = strings.TrimSpace(clientIP)
	if clientIP == "" {
		return ""
	}
	return unifyCountryCode(ipgeo.GetCountryCode(clientIP))
}

// unifyCountryCode 走 ipgeo 规范化 + constants/country 表(与登录入库同一套简码).
func unifyCountryCode(raw string) string {
	code := ipgeo.NormalizeCountryCode(raw)
	if code == "" {
		return ""
	}
	if c, ok := country.Get(code); ok {
		return c.Code
	}
	// Geo 有码但静态表暂无条目时仍回传规范化 ISO,避免漏报
	return code
}
