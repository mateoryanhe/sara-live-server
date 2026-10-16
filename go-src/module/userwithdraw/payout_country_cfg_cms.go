package userwithdraw

import (
	"context"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/haipaydto"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

var haiPayUserWithdrawPayoutContinentOrder = []string{
	country.HaiPayPayoutRegionNorthAmerica,
	country.HaiPayPayoutRegionEurope,
	country.HaiPayPayoutRegionSouthAmerica,
	country.HaiPayPayoutRegionAsia,
	country.HaiPayPayoutRegionMiddleEast,
	country.HaiPayPayoutRegionAfrica,
}

func GetHaiPayUserWithdrawPayoutCountryCfg(ctx context.Context, _ *haipaydto.GetUserWithdrawPayoutCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	_ = ctx
	groups := make(map[string]*haipaydto.CollectionCountryCfgGroup, len(haiPayUserWithdrawPayoutContinentOrder))
	for _, continent := range haiPayUserWithdrawPayoutContinentOrder {
		groups[continent] = &haipaydto.CollectionCountryCfgGroup{Continent: continent, Countries: []*haipaydto.CollectionCountryCfgItem{}}
	}
	configured := cfgdao.GetHaiPayUserWithdrawPayoutCountryCfgMapCached()
	for _, option := range country.ListHaiPayPayoutOptions() {
		continent := option.Region
		if continent == "" {
			continent = country.HaiPayPayoutRegion(option.CountryCode)
		}
		item := &haipaydto.CollectionCountryCfgItem{
			CountryCode:            option.CountryCode,
			Continent:              continent,
			SupportedCurrencies:    []string{option.Currency},
			CurrencyCode:           option.Currency,
			Enabled:                false,
			PaymentMethods:         make([]*haipaydto.CollectionPaymentMethodOption, 0),
			SelectedPaymentMethods: make([]*haipaydto.CollectionPaymentMethodSelection, 0, 1),
		}
		for _, method := range country.ListHaiPayPayoutMethods(option.Currency) {
			minAmount, maxAmount := "", ""
			if parts := strings.Split(strings.TrimSpace(method.Limit), "-"); len(parts) == 2 {
				minAmount, maxAmount = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			}
			item.PaymentMethods = append(item.PaymentMethods, &haipaydto.CollectionPaymentMethodOption{
				CurrencyCode: option.Currency,
				PayType:      method.AccountType,
				InBankCode:   method.BankCode,
				MinAmount:    minAmount,
				MaxAmount:    maxAmount,
				Description:  method.Description,
				Available:    true,
			})
		}
		if info, ok := country.Get(option.CountryCode); ok {
			item.CountryNameEn = info.NameEn
			item.CountryNameZh = info.NameZh
		}
		if stored := configured[option.CountryCode]; stored != nil && cfgdao.HaiPayUserWithdrawPayoutCountryCfgComplete(stored) {
			currencyCode := strings.ToUpper(strings.TrimSpace(stored.CurrencyCode))
			if currencyCode == option.Currency {
				item.ID = stored.CountryCode
				item.CurrencyCode = currencyCode
				item.Enabled = stored.Enabled
				accountType := strings.ToUpper(strings.TrimSpace(stored.AccountType))
				if _, ok := country.FindHaiPayPayoutMethod(currencyCode, accountType, stored.BankCode); ok {
					item.SelectedPaymentMethods = append(item.SelectedPaymentMethods, &haipaydto.CollectionPaymentMethodSelection{
						CurrencyCode: currencyCode,
						PayType:      accountType,
						InBankCode:   strings.TrimSpace(stored.BankCode),
					})
				}
			}
		}
		group := groups[continent]
		if group == nil {
			group = &haipaydto.CollectionCountryCfgGroup{Continent: continent, Countries: []*haipaydto.CollectionCountryCfgItem{}}
			groups[continent] = group
		}
		group.Countries = append(group.Countries, item)
	}
	continents := make([]*haipaydto.CollectionCountryCfgGroup, 0, len(groups))
	for _, continent := range haiPayUserWithdrawPayoutContinentOrder {
		if group := groups[continent]; group != nil && len(group.Countries) > 0 {
			continents = append(continents, group)
		}
	}
	return &haipaydto.GetCollectionCountryCfgRes{Continents: continents, PaymentTypes: []string{}}, nil
}

func SaveHaiPayUserWithdrawPayoutCountryCfg(ctx context.Context, req *haipaydto.SaveUserWithdrawPayoutCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	_ = ctx
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	currencyCode := strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
	accountType := strings.ToUpper(strings.TrimSpace(req.PayType))
	bankCode := strings.TrimSpace(req.InBankCode)
	if fixedType, ok := UserWithdrawPayoutFixedAccountType(countryCode); ok {
		if accountType != fixedType {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		accountType = fixedType
	}
	if !country.IsHaiPayPayoutCountry(countryCode) || country.HaiPayPayoutCurrency(countryCode) != currencyCode {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, ok := country.FindHaiPayPayoutMethod(currencyCode, accountType, bankCode); !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	appID, ok := country.LookupHaiPayAppID(currencyCode, cfgdao.HaiPayUseProdAppID())
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &rechargeentity.HaiPayUserWithdrawPayoutCountryCfg{
		CountryCode:  countryCode,
		CurrencyCode: currencyCode,
		AppId:        appID,
		AccountType:  accountType,
		BankCode:     bankCode,
		Enabled:      req.Enabled,
	}
	id, err := cfgdao.SaveHaiPayUserWithdrawPayoutCountryCfg(row)
	if err != nil {
		return nil, err
	}
	cfgdao.ReloadHaiPayUserWithdrawPayoutCountryCfgCache()
	return &haipaydto.SaveCollectionCountryCfgRes{Success: true, ID: id}, nil
}
