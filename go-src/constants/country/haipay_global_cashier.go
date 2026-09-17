package country

import "strings"

// haiPayGlobalCashierOptions 来自 HaiPay 全球收银台“地区编码”表及官网国家目录。
// CountryCode 使用 ISO 3166-1 alpha-2（EU 是 HaiPay 使用的地区特例）；
// Currencies 是 /global/cashier/collect/apply 接口允许该地区提交的交易币种。
var haiPayGlobalCashierOptions = []HaiPayCollectionOption{
	{CountryCode: "US", Currencies: []string{"USD"}},
	{CountryCode: "MX", Currencies: []string{"MXN"}},

	{CountryCode: "AT", Currencies: []string{"EUR"}},
	{CountryCode: "BE", Currencies: []string{"EUR"}},
	{CountryCode: "GB", Currencies: []string{"GBP"}},
	{CountryCode: "NL", Currencies: []string{"EUR"}},
	{CountryCode: "PL", Currencies: []string{"EUR"}},
	{CountryCode: "TR", Currencies: []string{"TRY"}},
	{CountryCode: "IT", Currencies: []string{"EUR"}},
	{CountryCode: "EU", Currencies: []string{"EUR"}},

	{CountryCode: "BR", Currencies: []string{"BRL"}},

	{CountryCode: "HK", Currencies: []string{"HKD", "USD"}},
	{CountryCode: "SG", Currencies: []string{"SGD", "USD"}},
	{CountryCode: "TW", Currencies: []string{"TWD", "USD"}},
	{CountryCode: "JP", Currencies: []string{"JPY", "USD"}},
	{CountryCode: "KR", Currencies: []string{"KRW", "USD"}},
	{CountryCode: "PH", Currencies: []string{"PHP"}},
	{CountryCode: "TH", Currencies: []string{"THB", "USD"}},
	{CountryCode: "VN", Currencies: []string{"VND", "USD"}},
	{CountryCode: "ID", Currencies: []string{"IDR", "USD"}},
	{CountryCode: "MY", Currencies: []string{"MYR", "USD"}},
	{CountryCode: "PK", Currencies: []string{"PKR"}},
	{CountryCode: "BD", Currencies: []string{"BDT"}},
	{CountryCode: "IN", Currencies: []string{"INR"}},

	{CountryCode: "EG", Currencies: []string{"EGP"}},
	{CountryCode: "SA", Currencies: []string{"SAR"}},
	{CountryCode: "KW", Currencies: []string{"KWD"}},
	{CountryCode: "BH", Currencies: []string{"BHD"}},
	{CountryCode: "AE", Currencies: []string{"AED"}},
	{CountryCode: "OM", Currencies: []string{"OMR"}},
	{CountryCode: "QA", Currencies: []string{"QAR"}},
	{CountryCode: "JO", Currencies: []string{"JOD"}},
	{CountryCode: "IQ", Currencies: []string{"IQD"}},

	{CountryCode: "NG", Currencies: []string{"NGN"}},
	{CountryCode: "GH", Currencies: []string{"GHS"}},
	{CountryCode: "KE", Currencies: []string{"KES"}},
	{CountryCode: "ZA", Currencies: []string{"ZAR"}},
	{CountryCode: "CM", Currencies: []string{"XAF"}},
	{CountryCode: "TZ", Currencies: []string{"TZS"}},
}

var haiPayGlobalCashierByCountry map[string]HaiPayCollectionOption
var haiPayGlobalCashierCountryByCurrency map[string]string

func init() {
	haiPayGlobalCashierByCountry = make(map[string]HaiPayCollectionOption, len(haiPayGlobalCashierOptions))
	haiPayGlobalCashierCountryByCurrency = make(map[string]string, len(haiPayGlobalCashierOptions))
	for i := range haiPayGlobalCashierOptions {
		item := &haiPayGlobalCashierOptions[i]
		item.CountryCode = normalizeCode(item.CountryCode)
		for currencyIndex := range item.Currencies {
			item.Currencies[currencyIndex] = strings.ToUpper(strings.TrimSpace(item.Currencies[currencyIndex]))
		}
		haiPayGlobalCashierByCountry[item.CountryCode] = *item
		for _, currency := range item.Currencies {
			if existing, exists := haiPayGlobalCashierCountryByCurrency[currency]; !exists {
				haiPayGlobalCashierCountryByCurrency[currency] = item.CountryCode
			} else if existing != item.CountryCode {
				haiPayGlobalCashierCountryByCurrency[currency] = ""
			}
		}
	}
	for currency, countryCode := range haiPayGlobalCashierCountryByCurrency {
		if countryCode == "" {
			delete(haiPayGlobalCashierCountryByCurrency, currency)
		}
	}
}

func ListHaiPayGlobalCashierOptions() []HaiPayCollectionOption {
	out := make([]HaiPayCollectionOption, 0, len(haiPayGlobalCashierOptions))
	for _, item := range haiPayGlobalCashierOptions {
		item.Currencies = append([]string(nil), item.Currencies...)
		out = append(out, item)
	}
	return out
}

func IsHaiPayGlobalCashierRegion(code string) bool {
	_, ok := haiPayGlobalCashierByCountry[normalizeCode(code)]
	return ok
}

// HaiPayGlobalCashierCurrency 返回默认交易币种；支持 USD 时沿用“优先 USD”的默认策略。
func HaiPayGlobalCashierCurrency(countryCode string) string {
	currencies := haiPayGlobalCashierByCountry[normalizeCode(countryCode)].Currencies
	for _, currency := range currencies {
		if currency == "USD" {
			return currency
		}
	}
	if len(currencies) == 0 {
		return ""
	}
	return currencies[0]
}

func HaiPayGlobalCashierCurrencies(countryCode string) []string {
	return append([]string(nil), haiPayGlobalCashierByCountry[normalizeCode(countryCode)].Currencies...)
}

// HaiPayGlobalCashierCountryFromCurrency 仅为唯一币种保留旧客户端兼容映射。
// EUR、USD 等共享币种必须由 App 上报国家/地区。
func HaiPayGlobalCashierCountryFromCurrency(currency string) string {
	return haiPayGlobalCashierCountryByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
}

func HaiPayGlobalCashierContinent(countryCode string) string {
	switch normalizeCode(countryCode) {
	case "US", "MX":
		return HaiPayPayoutRegionNorthAmerica
	case "AT", "BE", "GB", "NL", "PL", "TR", "IT", "EU":
		return HaiPayPayoutRegionEurope
	case "BR":
		return HaiPayPayoutRegionSouthAmerica
	case "HK", "SG", "TW", "JP", "KR", "PH", "TH", "VN", "ID", "MY", "PK", "BD", "IN":
		return HaiPayPayoutRegionAsia
	case "EG", "SA", "KW", "BH", "AE", "OM", "QA", "JO", "IQ":
		return HaiPayPayoutRegionMiddleEast
	case "NG", "GH", "KE", "ZA", "CM", "TZ":
		return HaiPayPayoutRegionAfrica
	default:
		return ""
	}
}
