package country

import "strings"

// HaiPayPayoutCountryCodes 支持 HaiPay 本地代付的国家/地区简码(有序).
// 对应接口路径 /{currency}/pay/apply；与收银台 region 白名单可不同.
var HaiPayPayoutCountryCodes = []string{
	"ID", "PH", "MY", "IN", "TH", "VN", "PK", "BR", "KR", "TW",
	"AE", "EG", "TR", "KE", "MX", "NG", "ZA", "BD", "TZ",
}

// countryCode -> ISO 4217 代付币种(小写路径用 ToLower)
var haiPayPayoutCurrencyByCountry = map[string]string{
	"ID": "IDR",
	"PH": "PHP",
	"MY": "MYR",
	"IN": "INR",
	"TH": "THB",
	"VN": "VND",
	"PK": "PKR",
	"BR": "BRL",
	"KR": "KRW",
	"TW": "TWD",
	"AE": "AED",
	"EG": "EGP",
	"TR": "TRY",
	"KE": "KES",
	"MX": "MXN",
	"NG": "NGN",
	"ZA": "ZAR",
	"BD": "BDT",
	"TZ": "TZS",
}

var haiPayPayoutCountryByCurrency map[string]string
var haiPayPayoutCountrySet map[string]struct{}

func init() {
	haiPayPayoutCountrySet = make(map[string]struct{}, len(HaiPayPayoutCountryCodes))
	haiPayPayoutCountryByCurrency = make(map[string]string, len(haiPayPayoutCurrencyByCountry))
	for _, code := range HaiPayPayoutCountryCodes {
		c := normalizeCode(code)
		haiPayPayoutCountrySet[c] = struct{}{}
		if cur, ok := haiPayPayoutCurrencyByCountry[c]; ok {
			haiPayPayoutCountryByCurrency[strings.ToUpper(cur)] = c
		}
	}
}

// IsHaiPayPayoutCountry 是否支持 HaiPay 代付的国家简码.
func IsHaiPayPayoutCountry(code string) bool {
	_, ok := haiPayPayoutCountrySet[normalizeCode(code)]
	return ok
}

// HaiPayPayoutCurrency 国家简码 → 代付币种(如 ID→IDR);不支持返回空.
func HaiPayPayoutCurrency(countryCode string) string {
	return haiPayPayoutCurrencyByCountry[normalizeCode(countryCode)]
}

// HaiPayPayoutCountryFromCurrency 代付币种 → 国家简码(如 IDR→ID);不支持返回空.
func HaiPayPayoutCountryFromCurrency(currency string) string {
	return haiPayPayoutCountryByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
}

// ListHaiPayPayoutCountries 返回可代付国家列表(含名称;缺表项跳过).
func ListHaiPayPayoutCountries() []Country {
	out := make([]Country, 0, len(HaiPayPayoutCountryCodes))
	for _, code := range HaiPayPayoutCountryCodes {
		if c, ok := Get(code); ok {
			out = append(out, c)
			continue
		}
		code = normalizeCode(code)
		out = append(out, Country{
			Code:     code,
			NameEn:   code,
			FlagIcon: FlagFileName(code),
		})
	}
	return out
}
