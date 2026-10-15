package userwithdraw

import (
	"context"

	"xr-game-server/constants/country"
	"xr-game-server/dto/userwithdrawdto"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

func AppListWithdrawCurrencies(_ context.Context, _ *userwithdrawdto.AppListWithdrawCurrenciesReq) (*userwithdrawdto.AppListWithdrawCurrenciesRes, error) {
	return &userwithdrawdto.AppListWithdrawCurrenciesRes{
		List: buildWithdrawCurrencyItems(),
	}, nil
}

func buildWithdrawCurrencyItems() []*userwithdrawdto.AppWithdrawCurrencyItem {
	options := listAppVisibleHaiPayPayoutOptions()
	version := countryflagdeploy.CurrentVersion()
	out := make([]*userwithdrawdto.AppWithdrawCurrencyItem, 0, len(options))
	for i, option := range options {
		c, ok := country.Get(option.CountryCode)
		if !ok {
			c = country.Country{Code: option.CountryCode, NameEn: option.CountryCode}
		}
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		out = append(out, &userwithdrawdto.AppWithdrawCurrencyItem{
			CountryCode: c.Code,
			Currency:    option.Currency,
			Continent:   option.Region,
			NameEn:      c.NameEn,
			NameZh:      c.NameZh,
			Icon:        icon,
			Sort:        i + 1,
		})
	}
	return out
}
