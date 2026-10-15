package userwithdraw

import (
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
)

func listAppVisibleHaiPayPayoutOptions() []country.HaiPayPayoutOption {
	configured := cfgdao.GetHaiPayUserWithdrawPayoutCountryCfgMapCached()
	out := make([]country.HaiPayPayoutOption, 0, len(configured))
	for _, option := range country.ListHaiPayPayoutOptions() {
		stored := configured[option.CountryCode]
		if stored == nil || !stored.Enabled || !cfgdao.HaiPayUserWithdrawPayoutCountryCfgComplete(stored) {
			continue
		}
		if strings.ToUpper(strings.TrimSpace(stored.CurrencyCode)) != option.Currency {
			continue
		}
		out = append(out, option)
	}
	return out
}
