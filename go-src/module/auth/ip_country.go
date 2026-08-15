package auth

import (
	"strings"

	"xr-game-server/entity/user"
	"xr-game-server/module/ipgeo"
)

func applyRegisterIpInfo(account *entity.Account, ip string) {
	if account == nil {
		return
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return
	}
	if account.RegisterIp == "" {
		account.SetRegisterIp(ip)
	}
	if account.RegisterCountry == "" {
		if country := ipgeo.GetCountryName(ip); country != "" {
			account.SetRegisterCountry(country)
		}
	}
	applyLoginIpInfo(account, ip)
}

func applyLoginIpInfo(account *entity.Account, ip string) {
	if account == nil {
		return
	}
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return
	}
	account.SetIp(ip)
	if country := ipgeo.GetCountryName(ip); country != "" {
		account.SetLoginCountry(country)
	}
}
