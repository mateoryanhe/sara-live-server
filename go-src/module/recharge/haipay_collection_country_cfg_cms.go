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

var haiPayCollectionContinentOrder = []string{
	country.HaiPayPayoutRegionNorthAmerica,
	country.HaiPayPayoutRegionEurope,
	country.HaiPayPayoutRegionSouthAmerica,
	country.HaiPayPayoutRegionAsia,
	country.HaiPayPayoutRegionMiddleEast,
	country.HaiPayPayoutRegionAfrica,
}

var haiPayCollectionPaymentTypes = []string{
	"CASHIER", "PAYMENT_GATEWAY", "BANK_TRANSFER", "BANK_ACCOUNT", "EWALLET", "QR", "VA",
}

// haiPayCollectionDefaultAppIDs 是币商本地代收各币种在 HaiPay 后台创建的业务 ID。
// 普通用户全球收银台 AppId 统一在 HaiPay 全局配置中维护。
var haiPayCollectionDefaultAppIDs = map[string]int64{
	"IDR":  25238,
	"PHP":  25239,
	"USD":  25240,
	"KRW":  25241,
	"EUR":  25242,
	"TWD":  25243,
	"VND":  25244,
	"THB":  25245,
	"SAR":  25246,
	"KWD":  25247,
	"OMR":  25248,
	"QAR":  25249,
	"AED":  25250,
	"BHD":  25251,
	"USDT": 25252,
	"BRL":  25253,
	"HKD":  25254,
	"INR":  25255,
	"MYR":  25256,
	"TRY":  25257,
	"PKR":  25258,
	"SGD":  25259,
	"RUB":  25260,
	"EGP":  25261,
	"JPY":  25262,
	"GBP":  25263,
	"PLN":  25264,
	"CAD":  25265,
	"JOD":  25266,
	"IQD":  25267,
	"BDT":  25268,
	"USDC": 25269,
	"MXN":  25270,
	"NGN":  25271,
}

func isHaiPayGlobalCashierBiz(bizType entity.HaiPayBizType) bool {
	return bizType == entity.HaiPayBizTypeNormalCollection
}

func haiPayCollectionOptionsForBiz(bizType entity.HaiPayBizType) []country.HaiPayCollectionOption {
	if isHaiPayGlobalCashierBiz(bizType) {
		return country.ListHaiPayGlobalCashierOptions()
	}
	return country.ListHaiPayCollectionOptions()
}

func haiPayCollectionContinentForBiz(bizType entity.HaiPayBizType, countryCode string) string {
	if isHaiPayGlobalCashierBiz(bizType) {
		return country.HaiPayGlobalCashierContinent(countryCode)
	}
	return country.HaiPayCollectionContinent(countryCode)
}

func haiPayCollectionDefaultCurrencyForBiz(bizType entity.HaiPayBizType, countryCode string) string {
	if isHaiPayGlobalCashierBiz(bizType) {
		return country.HaiPayGlobalCashierCurrency(countryCode)
	}
	return country.HaiPayCollectionCurrency(countryCode)
}

func haiPayCollectionCurrenciesForBiz(bizType entity.HaiPayBizType, countryCode string) []string {
	if isHaiPayGlobalCashierBiz(bizType) {
		return country.HaiPayGlobalCashierCurrencies(countryCode)
	}
	return country.HaiPayCollectionCurrencies(countryCode)
}

func isHaiPayCollectionRegionForBiz(bizType entity.HaiPayBizType, countryCode string) bool {
	if isHaiPayGlobalCashierBiz(bizType) {
		return country.IsHaiPayGlobalCashierRegion(countryCode)
	}
	return country.IsHaiPayRegion(countryCode)
}

func GetHaiPayCollectionCountryCfg(ctx context.Context, req *haipaydto.GetCollectionCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	return getHaiPayCollectionCountryCfg(ctx, req, entity.HaiPayBizTypeNormalCollection)
}

func getHaiPayCollectionCountryCfg(
	_ context.Context, _ *haipaydto.GetCollectionCountryCfgReq, bizType entity.HaiPayBizType,
) (*haipaydto.GetCollectionCountryCfgRes, error) {
	groups := make(map[string]*haipaydto.CollectionCountryCfgGroup, len(haiPayCollectionContinentOrder))
	for _, continent := range haiPayCollectionContinentOrder {
		groups[continent] = &haipaydto.CollectionCountryCfgGroup{Continent: continent, Countries: []*haipaydto.CollectionCountryCfgItem{}}
	}
	configured := cfgdao.GetHaiPayCollectionCountryCfgMapCached(bizType)
	for _, option := range haiPayCollectionOptionsForBiz(bizType) {
		continent := haiPayCollectionContinentForBiz(bizType, option.CountryCode)
		item := &haipaydto.CollectionCountryCfgItem{
			CountryCode: option.CountryCode, Continent: continent,
			SupportedCurrencies:    append([]string(nil), option.Currencies...),
			CurrencyCode:           haiPayCollectionDefaultCurrencyForBiz(bizType, option.CountryCode),
			Enabled:                false,
			PaymentMethods:         make([]*haipaydto.CollectionPaymentMethodOption, 0),
			SelectedPaymentMethods: make([]*haipaydto.CollectionPaymentMethodSelection, 0),
		}
		if info, ok := country.Get(option.CountryCode); ok {
			item.CountryNameEn = info.NameEn
			item.CountryNameZh = info.NameZh
		}
		if aggregate := configured[option.CountryCode]; aggregate != nil && aggregate.Config != nil {
			storedCurrency := strings.ToUpper(strings.TrimSpace(aggregate.Config.CurrencyCode))
			if containsHaiPayString(option.Currencies, storedCurrency) {
				item.ID = aggregate.Config.ID
				item.CurrencyCode = storedCurrency
				item.Enabled = aggregate.Config.Enabled
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

	return &haipaydto.GetCollectionCountryCfgRes{
		Continents: continents, PaymentTypes: []string{}, DefaultAppIds: map[string]int64{},
	}, nil
}

func SaveHaiPayCollectionCountryCfg(ctx context.Context, req *haipaydto.SaveCollectionCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	return saveHaiPayCollectionCountryCfg(ctx, req, entity.HaiPayBizTypeNormalCollection)
}

func saveHaiPayCollectionCountryCfg(
	_ context.Context, req *haipaydto.SaveCollectionCountryCfgReq, bizType entity.HaiPayBizType,
) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	currencyCode := strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
	if !isHaiPayCollectionRegionForBiz(bizType, countryCode) ||
		!containsHaiPayString(haiPayCollectionCurrenciesForBiz(bizType, countryCode), currencyCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.SaveSection != haipaydto.SaveCollectionCountryCfgSectionBasic || !isHaiPayGlobalCashierBiz(bizType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return saveHaiPayCollectionCountryBasicCfg(req, bizType, countryCode, currencyCode)
}

func saveHaiPayCollectionCountryBasicCfg(
	req *haipaydto.SaveCollectionCountryCfgReq, bizType entity.HaiPayBizType, countryCode, currencyCode string,
) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	config := &entity.HaiPayCollectionCountryCfg{
		BizType: bizType, CountryCode: countryCode, CurrencyCode: currencyCode, Enabled: req.Enabled,
	}
	id, err := cfgdao.SaveHaiPayCollectionCountryBasicCfg(config)
	if err != nil {
		return nil, err
	}
	cfgdao.ReloadHaiPayCollectionCountryCfgCache()
	return &haipaydto.SaveCollectionCountryCfgRes{Success: true, ID: id}, nil
}

func containsHaiPayString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
