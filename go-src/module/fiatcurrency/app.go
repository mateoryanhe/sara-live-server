package fiatcurrency

import (
	"context"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

// GetAppList App 查询普通用户 HaiPay 支付区域。
func GetAppList(_ context.Context, req *fiatcurrencydto.AppFiatCurrencyListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	if req != nil && req.TypeFilter == int(fiatentity.FiatCurrencyTypeCrypto) {
		return &fiatcurrencydto.AppFiatCurrencyListRes{List: []*fiatcurrencydto.AppFiatCurrencyItem{}}, nil
	}
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: buildHaiPayRegionItems(true, false)}, nil
}

// GetCoinMerchantAppList App 查询币商专用 HaiPay 支付区域。
func GetCoinMerchantAppList(_ context.Context, _ *fiatcurrencydto.CoinMerchantPaymentRegionListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: buildHaiPayRegionItems(true, true)}, nil
}

// GetHaiPayRegionList CMS 第三方充值测试使用普通用户的 HaiPay region 选型列表。
func GetHaiPayRegionList(_ context.Context, _ *fiatcurrencydto.HaiPayRegionListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: buildHaiPayRegionItems(false, false)}, nil
}

func buildHaiPayRegionItems(enabledOnly, coinMerchant bool) []*fiatcurrencydto.AppFiatCurrencyItem {
	options := country.ListHaiPayGlobalCashierOptions()
	configured := cfgdao.GetHaiPayCollectionCountryCfgMapCached(rechargeentity.HaiPayBizTypeNormalCollection)
	coinMerchantConfigured := map[string]*rechargeentity.HaiPayCoinMerchantCollectionCfg(nil)
	if coinMerchant {
		options = country.ListHaiPayCollectionOptions()
		configured = nil
		coinMerchantConfigured = cfgdao.GetHaiPayCoinMerchantCollectionCfgMapCached()
	}
	version := countryflagdeploy.CurrentVersion()
	list := make([]*fiatcurrencydto.AppFiatCurrencyItem, 0, len(options))
	for i, option := range options {
		c, ok := country.Get(option.CountryCode)
		if !ok {
			c = country.Country{Code: option.CountryCode, NameEn: option.CountryCode}
		}
		icon := ""
		if rel := country.RelPath(option.CountryCode, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		continent := country.HaiPayGlobalCashierContinent(option.CountryCode)
		fallbackCurrency := country.HaiPayGlobalCashierCurrency(option.CountryCode)
		if coinMerchant {
			continent = country.HaiPayCollectionContinent(option.CountryCode)
			fallbackCurrency = country.HaiPayCollectionCurrency(option.CountryCode)
		}
		item := &fiatcurrencydto.AppFiatCurrencyItem{
			CurrencyCode: c.Code,
			Continent:    continent,
			Name:         c.NameEn,
			NameEn:       c.NameEn,
			NameZh:       c.NameZh,
			Symbol:       c.Code,
			Icon:         icon,
			CurrencyType: uint8(fiatentity.FiatCurrencyTypeFiat),
			Sort:         i + 1,
		}
		aggregate := configured[option.CountryCode]
		coinMerchantConfig := coinMerchantConfigured[option.CountryCode]
		if enabledOnly {
			if coinMerchant {
				_, _, enabled, valid := resolveConfiguredCollectionMethods(
					option.CountryCode, fallbackCurrency, coinMerchantConfig,
				)
				if !valid || !enabled {
					continue
				}
			} else if aggregate == nil || aggregate.Config == nil || !aggregate.Config.Enabled {
				continue
			}
			// App 只需要区域。币种、AppId 与支付方式在建单时按业务类型和地区从缓存读取。
			list = append(list, item)
			continue
		}

		if !coinMerchant {
			currencyCode := fallbackCurrency
			if aggregate != nil && aggregate.Config != nil {
				storedCurrency := strings.ToUpper(strings.TrimSpace(aggregate.Config.CurrencyCode))
				for _, supportedCurrency := range option.Currencies {
					if storedCurrency == supportedCurrency {
						currencyCode = storedCurrency
						break
					}
				}
			}
			item.FiatCurrencyCode = currencyCode
			item.FiatCurrencyCodes = append([]string(nil), option.Currencies...)
			list = append(list, item)
			continue
		}

		currencyCode, paymentMethods, _, valid := resolveConfiguredCollectionMethods(
			option.CountryCode, fallbackCurrency, coinMerchantConfig,
		)
		if !valid {
			currencyCode = fallbackCurrency
			paymentMethods = nil
		}
		payTypes := make([]string, 0, len(paymentMethods))
		inBankCodes := make([]string, 0, len(paymentMethods))
		seenTypes := make(map[string]struct{})
		for _, method := range paymentMethods {
			if method == nil {
				continue
			}
			if _, ok := seenTypes[method.PayType]; !ok {
				seenTypes[method.PayType] = struct{}{}
				payTypes = append(payTypes, method.PayType)
			}
			inBankCodes = append(inBankCodes, method.InBankCode)
		}
		item.FiatCurrencyCode = currencyCode
		item.FiatCurrencyCodes = append([]string(nil), option.Currencies...)
		item.PayType = strings.Join(payTypes, ",")
		item.InBankCode = strings.Join(inBankCodes, ",")
		item.PayTypes = payTypes
		item.InBankCodes = inBankCodes
		item.PaymentMethods = paymentMethods
		list = append(list, item)
	}
	return list
}

func resolveConfiguredCollectionMethods(
	countryCode, fallbackCurrency string, config *rechargeentity.HaiPayCoinMerchantCollectionCfg,
) (currency string, paymentMethods []*fiatcurrencydto.AppHaiPayPaymentMethod, enabled, valid bool) {
	if !cfgdao.HaiPayCoinMerchantCollectionCfgComplete(config) {
		return fallbackCurrency, nil, false, false
	}
	currency = strings.ToUpper(strings.TrimSpace(config.CurrencyCode))
	enabled = config.Enabled
	payType := strings.ToUpper(strings.TrimSpace(config.PayType))
	canonicalCode, ok := country.ResolveHaiPayCollectionPaymentMethodCode(
		countryCode, currency, payType, config.InBankCode,
	)
	if !ok || canonicalCode == "" {
		return currency, nil, enabled, false
	}
	item := &fiatcurrencydto.AppHaiPayPaymentMethod{PayType: payType, InBankCode: canonicalCode}
	catalog := country.ListHaiPayCollectionPaymentMethods(countryCode)
	for _, option := range catalog {
		if option.Available && option.CurrencyCode == currency && option.PayType == payType && strings.EqualFold(option.InBankCode, canonicalCode) {
			item.MinAmount = option.MinAmount
			item.MaxAmount = option.MaxAmount
			item.Description = option.Description
			paymentMethods = append(paymentMethods, item)
			break
		}
	}
	if len(paymentMethods) != 1 {
		return currency, nil, enabled, false
	}
	return currency, paymentMethods, enabled, true
}
