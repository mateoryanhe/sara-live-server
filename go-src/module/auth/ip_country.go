package auth

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/entity/user"
	"xr-game-server/module/ipgeo"
)

func applyRegisterIpInfo(account *entity.Account, r *ghttp.Request) {
	if account == nil || r == nil {
		return
	}
	location := userloginlocationdao.GetByUserId(account.ID)
	if location == nil {
		return
	}
	ip := strings.TrimSpace(r.GetClientIp())
	if ip == "" {
		return
	}
	if location.RegisterIp == "" {
		location.SetRegisterIP(ip)
	}
	if location.RegisterCountry == "" {
		if code := resolveRequestCountryCode(r, ip); code != "" {
			location.SetRegisterCountry(code)
		}
	}
	applyLoginLocation(location, r, ip)
	userloginlocationdao.Publish(location)
}

func applyLoginIpInfo(account *entity.Account, r *ghttp.Request) {
	if account == nil || r == nil {
		return
	}
	location := userloginlocationdao.GetByUserId(account.ID)
	if location == nil {
		return
	}
	ip := strings.TrimSpace(r.GetClientIp())
	if ip == "" {
		return
	}
	applyLoginLocation(location, r, ip)
	userloginlocationdao.Publish(location)
}

func applyLoginLocation(location *entity.UserLoginLocation, r *ghttp.Request, ip string) {
	location.SetIP(ip)
	if code := resolveRequestCountryCode(r, ip); code != "" {
		location.SetLoginCountry(code)
	}
}

func resolveRequestCountryCode(r *ghttp.Request, ip string) string {
	cfCountry := ""
	if r != nil {
		// Cloudflare: CF-IPCountry,如 US;HTTP 头大小写不敏感
		cfCountry = r.Header.Get("CF-IPCountry")
	}
	return ipgeo.ResolveCountryCode(ip, cfCountry)
}
