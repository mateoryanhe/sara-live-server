package recharge

import (
	"context"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/haipaydto"
	entity "xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

func GetHaiPayCoinMerchantCollectionCountryCfg(ctx context.Context, req *haipaydto.GetCollectionCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	_ = ctx
	_ = req
	groups := make(map[string]*haipaydto.CollectionCountryCfgGroup, len(haiPayCollectionContinentOrder))
	for _, continent := range haiPayCollectionContinentOrder {
		groups[continent] = &haipaydto.CollectionCountryCfgGroup{Continent: continent, Countries: []*haipaydto.CollectionCountryCfgItem{}}
	}
	configured := cfgdao.GetHaiPayCoinMerchantCollectionCfgMapCached()
	for _, option := range country.ListHaiPayCollectionOptions() {
		continent := country.HaiPayCollectionContinent(option.CountryCode)
		defaultCurrency := country.HaiPayCollectionCurrency(option.CountryCode)
		item := &haipaydto.CollectionCountryCfgItem{
			CountryCode: option.CountryCode, Continent: continent,
			SupportedCurrencies:    append([]string(nil), option.Currencies...),
			CurrencyCode:           defaultCurrency,
			AppId:                  haiPayCollectionDefaultAppIDs[defaultCurrency],
			Enabled:                false,
			PaymentMethods:         make([]*haipaydto.CollectionPaymentMethodOption, 0),
			SelectedPaymentMethods: make([]*haipaydto.CollectionPaymentMethodSelection, 0, 1),
		}
		for _, method := range country.ListHaiPayCollectionPaymentMethods(option.CountryCode) {
			item.PaymentMethods = append(item.PaymentMethods, &haipaydto.CollectionPaymentMethodOption{
				CurrencyCode: method.CurrencyCode, PayType: method.PayType, InBankCode: method.InBankCode,
				MinAmount: method.MinAmount, MaxAmount: method.MaxAmount,
				Description: method.Description, Available: method.Available,
			})
		}
		if info, ok := country.Get(option.CountryCode); ok {
			item.CountryNameEn = info.NameEn
			item.CountryNameZh = info.NameZh
		}
		if stored := configured[option.CountryCode]; stored != nil {
			currencyCode := strings.ToUpper(strings.TrimSpace(stored.CurrencyCode))
			if containsHaiPayString(option.Currencies, currencyCode) {
				item.ID = stored.CountryCode
				item.CurrencyCode = currencyCode
				item.AppId = stored.AppId
				item.Enabled = stored.Enabled
				payType := strings.ToUpper(strings.TrimSpace(stored.PayType))
				if code, valid := country.ResolveHaiPayCollectionPaymentMethodCode(
					option.CountryCode, currencyCode, payType, stored.InBankCode,
				); valid && code != "" {
					item.SelectedPaymentMethods = append(item.SelectedPaymentMethods, &haipaydto.CollectionPaymentMethodSelection{
						CurrencyCode: currencyCode, PayType: payType, InBankCode: code,
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
	for _, continent := range haiPayCollectionContinentOrder {
		if group := groups[continent]; group != nil && len(group.Countries) > 0 {
			continents = append(continents, group)
		}
	}
	defaultAppIDs := make(map[string]int64, len(haiPayCollectionDefaultAppIDs))
	for currencyCode, appID := range haiPayCollectionDefaultAppIDs {
		defaultAppIDs[currencyCode] = appID
	}
	return &haipaydto.GetCollectionCountryCfgRes{
		Continents:    continents,
		PaymentTypes:  append([]string(nil), haiPayCollectionPaymentTypes...),
		DefaultAppIds: defaultAppIDs,
	}, nil
}

func SaveHaiPayCoinMerchantCollectionCountryCfg(ctx context.Context, req *haipaydto.SaveCoinMerchantCollectionCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	_ = ctx
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	currencyCode := strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
	payType := strings.ToUpper(strings.TrimSpace(req.PayType))
	if !country.IsHaiPayRegion(countryCode) ||
		!containsHaiPayString(country.HaiPayCollectionCurrencies(countryCode), currencyCode) || req.AppId <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	canonicalCode, ok := country.ResolveHaiPayCollectionPaymentMethodCode(
		countryCode, currencyCode, payType, req.InBankCode,
	)
	if !ok || canonicalCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.HaiPayCoinMerchantCollectionCfg{
		CountryCode: countryCode, CurrencyCode: currencyCode, AppId: req.AppId,
		PayType: payType, InBankCode: canonicalCode, Enabled: req.Enabled,
	}
	id, err := cfgdao.SaveHaiPayCoinMerchantCollectionCfg(row)
	if err != nil {
		return nil, err
	}
	cfgdao.ReloadHaiPayCoinMerchantCollectionCfgCache()
	return &haipaydto.SaveCollectionCountryCfgRes{Success: true, ID: id}, nil
}
