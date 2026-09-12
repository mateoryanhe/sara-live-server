package auth

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/entity/user"
	"xr-game-server/module/ipgeo"
)

func applyRegisterIpInfo(account *entity.Account, r *ghttp.Request) {
	if account == nil || r == nil {
		return
	}
	ip := strings.TrimSpace(r.GetClientIp())
	if ip == "" {
		return
	}
	if account.RegisterIp == "" {
		account.SetRegisterIp(ip)
	}
	if account.RegisterCountry == "" {
		if code := resolveRequestCountryCode(r, ip); code != "" {
			account.SetRegisterCountry(code)
		}
	}
	applyLoginIpInfo(account, r)
}

func applyLoginIpInfo(account *entity.Account, r *ghttp.Request) {
	if account == nil || r == nil {
		return
	}
	ip := strings.TrimSpace(r.GetClientIp())
	if ip == "" {
		return
	}
	account.SetIp(ip)
	if code := resolveRequestCountryCode(r, ip); code != "" {
		account.SetLoginCountry(code)
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
