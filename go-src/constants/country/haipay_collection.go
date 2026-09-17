package country

import "strings"

// HaiPayCollectionOption 描述 HaiPay 国家本地代收接口的国家/地区与支持币种。
// CountryCode 是 HaiPay region；Currencies 是官网列出的 ISO 4217 币种代码。
type HaiPayCollectionOption struct {
	CountryCode string
	Currencies  []string
}

// haiPayCollectionOptions 是 HaiPay 本地代收国家与币种的唯一数据源。
// 顺序按官网国家目录的洲分组排列；不包含只有代付接口的地区。
var haiPayCollectionOptions = []HaiPayCollectionOption{
	{CountryCode: "AT", Currencies: []string{"EUR", "USD"}},
	{CountryCode: "BE", Currencies: []string{"EUR", "USD"}},
	{CountryCode: "IT", Currencies: []string{"EUR", "USD"}},
	{CountryCode: "NL", Currencies: []string{"EUR", "USD"}},
	{CountryCode: "PL", Currencies: []string{"USD"}},
	{CountryCode: "TR", Currencies: []string{"TRY"}},

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
	{CountryCode: "IN", Currencies: []string{"INR", "USD"}},
	{CountryCode: "MY", Currencies: []string{"MYR", "USD"}},
	{CountryCode: "PK", Currencies: []string{"PKR", "USD"}},
	{CountryCode: "BD", Currencies: []string{"BDT"}},

	{CountryCode: "EG", Currencies: []string{"EGP", "USD"}},
	{CountryCode: "SA", Currencies: []string{"SAR", "USD"}},
	{CountryCode: "AE", Currencies: []string{"AED", "USD"}},
	{CountryCode: "KW", Currencies: []string{"KWD", "USD"}},
	{CountryCode: "QA", Currencies: []string{"QAR", "USD"}},
	{CountryCode: "OM", Currencies: []string{"OMR", "USD"}},
	{CountryCode: "BH", Currencies: []string{"BHD", "USD"}},
	{CountryCode: "JO", Currencies: []string{"JOD"}},
	{CountryCode: "IQ", Currencies: []string{"IQD"}},
}

// HaiPayRegionCodes 保留给旧调用方使用；数据由 haiPayCollectionOptions 自动生成。
var HaiPayRegionCodes []string

var haiPayCollectionByCountry map[string]HaiPayCollectionOption
var haiPayCollectionCountryByCurrency map[string]string

func init() {
	HaiPayRegionCodes = make([]string, 0, len(haiPayCollectionOptions))
	haiPayCollectionByCountry = make(map[string]HaiPayCollectionOption, len(haiPayCollectionOptions))
	haiPayCollectionCountryByCurrency = make(map[string]string, len(haiPayCollectionOptions))
	for i := range haiPayCollectionOptions {
		item := &haiPayCollectionOptions[i]
		item.CountryCode = normalizeCode(item.CountryCode)
		for currencyIndex := range item.Currencies {
			item.Currencies[currencyIndex] = strings.ToUpper(strings.TrimSpace(item.Currencies[currencyIndex]))
		}
		HaiPayRegionCodes = append(HaiPayRegionCodes, item.CountryCode)
		haiPayCollectionByCountry[item.CountryCode] = *item
		for _, currency := range item.Currencies {
			if existing, exists := haiPayCollectionCountryByCurrency[currency]; !exists {
				haiPayCollectionCountryByCurrency[currency] = item.CountryCode
			} else if existing != item.CountryCode {
				// EUR/USD 等共享币种不能反推出国家；App 必须上报地区。
				haiPayCollectionCountryByCurrency[currency] = ""
			}
		}
	}
	for currency, countryCode := range haiPayCollectionCountryByCurrency {
		if countryCode == "" {
			delete(haiPayCollectionCountryByCurrency, currency)
		}
	}
}

// HaiPayCollectionContinent 返回与 HaiPay 后台一致的地区分组。
func HaiPayCollectionContinent(countryCode string) string {
	switch normalizeCode(countryCode) {
	case "AT", "BE", "IT", "NL", "PL", "TR":
		return HaiPayPayoutRegionEurope
	case "BR":
		return HaiPayPayoutRegionSouthAmerica
	case "HK", "SG", "TW", "JP", "KR", "PH", "TH", "VN", "ID", "IN", "MY", "PK", "BD":
		return HaiPayPayoutRegionAsia
	case "EG", "SA", "AE", "KW", "QA", "OM", "BH", "JO", "IQ":
		return HaiPayPayoutRegionMiddleEast
	default:
		return ""
	}
}

// ListHaiPayCollectionOptions 返回代收国家和币种目录的深拷贝。
func ListHaiPayCollectionOptions() []HaiPayCollectionOption {
	out := make([]HaiPayCollectionOption, 0, len(haiPayCollectionOptions))
	for _, item := range haiPayCollectionOptions {
		item.Currencies = append([]string(nil), item.Currencies...)
		out = append(out, item)
	}
	return out
}

// IsHaiPayRegion 是否为 HaiPay 收银台合法 region 简码。
func IsHaiPayRegion(code string) bool {
	_, ok := haiPayCollectionByCountry[normalizeCode(code)]
	return ok
}

// HaiPayCollectionCurrency 返回国家/地区的下单币种：支持 USD 时优先 USD，否则使用官网列表中的当地币。
func HaiPayCollectionCurrency(countryCode string) string {
	currencies := haiPayCollectionByCountry[normalizeCode(countryCode)].Currencies
	if len(currencies) == 0 {
		return ""
	}
	for _, currency := range currencies {
		if currency == "USD" {
			return currency
		}
	}
	return currencies[0]
}

// HaiPayCollectionCurrencies 返回国家/地区支持的全部代收币种。
func HaiPayCollectionCurrencies(countryCode string) []string {
	return append([]string(nil), haiPayCollectionByCountry[normalizeCode(countryCode)].Currencies...)
}

// HaiPayCollectionCountryFromCurrency 仅在币种唯一对应一个地区时返回 HaiPay region。
// USD、EUR 等共享币种必须由 App 上报国家/地区，不能在服务端猜测。
func HaiPayCollectionCountryFromCurrency(currency string) string {
	return haiPayCollectionCountryByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
}

// ListHaiPayRegions 返回 HaiPay 可选国家/地区，顺序与代收目录一致。
func ListHaiPayRegions() []Country {
	out := make([]Country, 0, len(haiPayCollectionOptions))
	for _, item := range haiPayCollectionOptions {
		if c, ok := Get(item.CountryCode); ok {
			out = append(out, c)
			continue
		}
		out = append(out, Country{
			Code:     item.CountryCode,
			NameEn:   item.CountryCode,
			FlagIcon: FlagFileName(item.CountryCode),
		})
	}
	return out
}

// NormalizeHaiPayRegionHint 将 region 规范化；非法返回空。
func NormalizeHaiPayRegionHint(code string) string {
	code = normalizeCode(code)
	if !IsHaiPayRegion(code) {
		return ""
	}
	return code
}
