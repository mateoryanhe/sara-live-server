package country

import "strings"

// HaiPayCollectionOption 描述全球收银台的国家/地区与支持币种。
// CountryCode 是 HaiPay region；Currencies 是官网列出的 ISO 4217 币种代码。
type HaiPayCollectionOption struct {
	CountryCode string
	Currencies  []string
}

// haiPayCollectionOptions 是 HaiPay 代收国家与币种的唯一数据源。
// 顺序与全球收银台“地区编码”表一致。
var haiPayCollectionOptions = []HaiPayCollectionOption{
	{CountryCode: "AE", Currencies: []string{"AED"}},
	{CountryCode: "AT", Currencies: []string{"EUR"}},
	{CountryCode: "BE", Currencies: []string{"EUR"}},
	{CountryCode: "BH", Currencies: []string{"BHD"}},
	{CountryCode: "BR", Currencies: []string{"BRL"}},
	{CountryCode: "EG", Currencies: []string{"EGP"}},
	{CountryCode: "GB", Currencies: []string{"GBP"}},
	{CountryCode: "HK", Currencies: []string{"HKD", "USD"}},
	{CountryCode: "ID", Currencies: []string{"IDR", "USD"}},
	{CountryCode: "IN", Currencies: []string{"INR"}},
	{CountryCode: "JP", Currencies: []string{"JPY", "USD"}},
	{CountryCode: "KR", Currencies: []string{"KRW", "USD"}},
	{CountryCode: "KW", Currencies: []string{"KWD"}},
	{CountryCode: "MY", Currencies: []string{"MYR", "USD"}},
	{CountryCode: "NL", Currencies: []string{"EUR"}},
	{CountryCode: "OM", Currencies: []string{"OMR"}},
	{CountryCode: "PH", Currencies: []string{"PHP"}},
	{CountryCode: "PK", Currencies: []string{"PKR"}},
	{CountryCode: "PL", Currencies: []string{"EUR"}},
	{CountryCode: "QA", Currencies: []string{"QAR"}},
	{CountryCode: "SA", Currencies: []string{"SAR"}},
	{CountryCode: "SG", Currencies: []string{"SGD", "USD"}},
	{CountryCode: "TH", Currencies: []string{"THB", "USD"}},
	{CountryCode: "TR", Currencies: []string{"TRY"}},
	{CountryCode: "TW", Currencies: []string{"TWD", "USD"}},
	{CountryCode: "US", Currencies: []string{"USD"}},
	{CountryCode: "VN", Currencies: []string{"VND", "USD"}},
	{CountryCode: "IT", Currencies: []string{"EUR"}},
	{CountryCode: "EU", Currencies: []string{"EUR"}},
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
			if _, exists := haiPayCollectionCountryByCurrency[currency]; !exists {
				haiPayCollectionCountryByCurrency[currency] = item.CountryCode
			}
		}
	}
	// EUR/USD 对应多个地区；按区域不明确的旧币种入参分别回落 EU/US。
	haiPayCollectionCountryByCurrency["EUR"] = "EU"
	haiPayCollectionCountryByCurrency["USD"] = "US"
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

// HaiPayCollectionCountryFromCurrency 返回法币默认对应的 HaiPay region。
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
