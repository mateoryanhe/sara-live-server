package country

import "strings"

const (
	HaiPayPayoutAccountTypeBank    = "BANK_ACCOUNT"
	HaiPayPayoutAccountTypeEWallet = "EWALLET"

	HaiPayPayoutRegionNorthAmerica = "NORTH_AMERICA"
	HaiPayPayoutRegionEurope       = "EUROPE"
	HaiPayPayoutRegionSouthAmerica = "SOUTH_AMERICA"
	HaiPayPayoutRegionAsia         = "ASIA"
	HaiPayPayoutRegionMiddleEast   = "MIDDLE_EAST"
	HaiPayPayoutRegionAfrica       = "AFRICA"
)

// HaiPayPayoutWallet 描述已接入电子钱包的支付编码。
type HaiPayPayoutWallet struct {
	Code string
	Name string
}

// HaiPayPayoutOption 描述一个国家/地区当前已接入的代付能力。
type HaiPayPayoutOption struct {
	CountryCode  string
	Currency     string
	Region       string
	AccountTypes []string
	Wallets      []HaiPayPayoutWallet
}

// haiPayPayoutOptions 是代付国家、币种、地区和账户类型的唯一数据源。
// 历史币种若当前文档未列出电子钱包，保守地只开放银行账户。
// 美国 ACH 还需要路由编码和地址字段，当前只开放现有资料表能完整支持的电子钱包。
var haiPayPayoutOptions = []HaiPayPayoutOption{
	{
		CountryCode: "US", Currency: "USD", Region: HaiPayPayoutRegionNorthAmerica,
		AccountTypes: []string{HaiPayPayoutAccountTypeEWallet},
		Wallets: []HaiPayPayoutWallet{
			{Code: "VENMO", Name: "Venmo"},
			{Code: "ECASHAPP", Name: "CashApp"},
		},
	},
	{CountryCode: "MX", Currency: "MXN", Region: HaiPayPayoutRegionNorthAmerica, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "TR", Currency: "TRY", Region: HaiPayPayoutRegionEurope, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "BR", Currency: "BRL", Region: HaiPayPayoutRegionSouthAmerica, AccountTypes: []string{HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "ID", Currency: "IDR", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "PH", Currency: "PHP", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "MY", Currency: "MYR", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "IN", Currency: "INR", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "TH", Currency: "THB", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "VN", Currency: "VND", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "PK", Currency: "PKR", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "KR", Currency: "KRW", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "TW", Currency: "TWD", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "BD", Currency: "BDT", Region: HaiPayPayoutRegionAsia, AccountTypes: []string{HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "AE", Currency: "AED", Region: HaiPayPayoutRegionMiddleEast, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "EG", Currency: "EGP", Region: HaiPayPayoutRegionMiddleEast, AccountTypes: []string{HaiPayPayoutAccountTypeBank, HaiPayPayoutAccountTypeEWallet}},
	{CountryCode: "KE", Currency: "KES", Region: HaiPayPayoutRegionAfrica, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "NG", Currency: "NGN", Region: HaiPayPayoutRegionAfrica, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "ZA", Currency: "ZAR", Region: HaiPayPayoutRegionAfrica, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
	{CountryCode: "TZ", Currency: "TZS", Region: HaiPayPayoutRegionAfrica, AccountTypes: []string{HaiPayPayoutAccountTypeBank}},
}

// HaiPayPayoutCountryCodes 保留给旧调用方使用；数据由 haiPayPayoutOptions 自动生成。
var HaiPayPayoutCountryCodes []string

var haiPayPayoutByCountry map[string]HaiPayPayoutOption
var haiPayPayoutByCurrency map[string]HaiPayPayoutOption

func init() {
	HaiPayPayoutCountryCodes = make([]string, 0, len(haiPayPayoutOptions))
	haiPayPayoutByCountry = make(map[string]HaiPayPayoutOption, len(haiPayPayoutOptions))
	haiPayPayoutByCurrency = make(map[string]HaiPayPayoutOption, len(haiPayPayoutOptions))
	for _, item := range haiPayPayoutOptions {
		item.CountryCode = normalizeCode(item.CountryCode)
		item.Currency = strings.ToUpper(strings.TrimSpace(item.Currency))
		HaiPayPayoutCountryCodes = append(HaiPayPayoutCountryCodes, item.CountryCode)
		haiPayPayoutByCountry[item.CountryCode] = item
		haiPayPayoutByCurrency[item.Currency] = item
	}
}

func cloneHaiPayPayoutOption(item HaiPayPayoutOption) HaiPayPayoutOption {
	item.AccountTypes = append([]string(nil), item.AccountTypes...)
	item.Wallets = append([]HaiPayPayoutWallet(nil), item.Wallets...)
	return item
}

// ListHaiPayPayoutOptions 返回完整代付目录的副本。
func ListHaiPayPayoutOptions() []HaiPayPayoutOption {
	out := make([]HaiPayPayoutOption, 0, len(haiPayPayoutOptions))
	for _, item := range haiPayPayoutOptions {
		out = append(out, cloneHaiPayPayoutOption(item))
	}
	return out
}

func IsHaiPayPayoutCountry(code string) bool {
	_, ok := haiPayPayoutByCountry[normalizeCode(code)]
	return ok
}

func HaiPayPayoutCurrency(countryCode string) string {
	return haiPayPayoutByCountry[normalizeCode(countryCode)].Currency
}

func HaiPayPayoutCountryFromCurrency(currency string) string {
	return haiPayPayoutByCurrency[strings.ToUpper(strings.TrimSpace(currency))].CountryCode
}

func HaiPayPayoutRegion(countryCode string) string {
	return haiPayPayoutByCountry[normalizeCode(countryCode)].Region
}

func ListHaiPayPayoutAccountTypes(currency string) []string {
	item := haiPayPayoutByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
	return append([]string(nil), item.AccountTypes...)
}

// ListHaiPayPayoutWallets 返回币种当前已接入的钱包及其支付编码。
func ListHaiPayPayoutWallets(currency string) []HaiPayPayoutWallet {
	item := haiPayPayoutByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
	return append([]HaiPayPayoutWallet(nil), item.Wallets...)
}

func IsHaiPayPayoutAccountType(currency, accountType string) bool {
	accountType = strings.ToUpper(strings.TrimSpace(accountType))
	for _, item := range ListHaiPayPayoutAccountTypes(currency) {
		if item == accountType {
			return true
		}
	}
	return false
}

func ListHaiPayPayoutCountries() []Country {
	out := make([]Country, 0, len(haiPayPayoutOptions))
	for _, item := range haiPayPayoutOptions {
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
