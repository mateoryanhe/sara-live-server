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

// HaiPayPayoutWallet 描述电子钱包支付编码，保留给现有调用方使用。
type HaiPayPayoutWallet struct {
	Code string
	Name string
}

// HaiPayPayoutMethod 描述 HaiPay 当前可用的一条代付方式。
// AccountType 与 BankCode 必须成对使用，不能跨币种或账户类型混用。
type HaiPayPayoutMethod struct {
	AccountType string
	BankCode    string
	Limit       string
	Description string
}

// HaiPayPayoutOption 描述一个国家/地区当前已接入的代付能力。
type HaiPayPayoutOption struct {
	CountryCode  string
	Currency     string
	Region       string
	AccountTypes []string
	Wallets      []HaiPayPayoutWallet
	Methods      []HaiPayPayoutMethod
}

// HaiPayPayoutCountryCodes 保留给旧调用方使用；数据由 haiPayPayoutOptions 自动生成。
var HaiPayPayoutCountryCodes []string

var haiPayPayoutByCountry map[string]HaiPayPayoutOption
var haiPayPayoutByCurrency map[string]HaiPayPayoutOption

func init() {
	HaiPayPayoutCountryCodes = make([]string, 0, len(haiPayPayoutOptions))
	haiPayPayoutByCountry = make(map[string]HaiPayPayoutOption, len(haiPayPayoutOptions))
	haiPayPayoutByCurrency = make(map[string]HaiPayPayoutOption, len(haiPayPayoutOptions))
	for itemIndex := range haiPayPayoutOptions {
		item := &haiPayPayoutOptions[itemIndex]
		item.CountryCode = normalizeCode(item.CountryCode)
		item.Currency = strings.ToUpper(strings.TrimSpace(item.Currency))
		accountTypeSeen := make(map[string]struct{})
		item.AccountTypes = item.AccountTypes[:0]
		item.Wallets = item.Wallets[:0]
		for methodIndex := range item.Methods {
			method := &item.Methods[methodIndex]
			method.AccountType = strings.ToUpper(strings.TrimSpace(method.AccountType))
			method.BankCode = strings.TrimSpace(method.BankCode)
			method.Limit = strings.TrimSpace(method.Limit)
			method.Description = strings.TrimSpace(method.Description)
			if _, ok := accountTypeSeen[method.AccountType]; !ok {
				accountTypeSeen[method.AccountType] = struct{}{}
				item.AccountTypes = append(item.AccountTypes, method.AccountType)
			}
			if method.AccountType == HaiPayPayoutAccountTypeEWallet {
				item.Wallets = append(item.Wallets, HaiPayPayoutWallet{
					Code: method.BankCode,
					Name: method.Description,
				})
			}
		}
		HaiPayPayoutCountryCodes = append(HaiPayPayoutCountryCodes, item.CountryCode)
		haiPayPayoutByCountry[item.CountryCode] = *item
		haiPayPayoutByCurrency[item.Currency] = *item
	}
}

func cloneHaiPayPayoutOption(item HaiPayPayoutOption) HaiPayPayoutOption {
	item.AccountTypes = append([]string(nil), item.AccountTypes...)
	item.Wallets = append([]HaiPayPayoutWallet(nil), item.Wallets...)
	item.Methods = append([]HaiPayPayoutMethod(nil), item.Methods...)
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

// ListHaiPayPayoutMethods 返回币种当前可用的具体代付方式。
func ListHaiPayPayoutMethods(currency string) []HaiPayPayoutMethod {
	item := haiPayPayoutByCurrency[strings.ToUpper(strings.TrimSpace(currency))]
	return append([]HaiPayPayoutMethod(nil), item.Methods...)
}

// FindHaiPayPayoutMethod 校验并返回 accountType 与 bankCode 对应的官方可用方式。
func FindHaiPayPayoutMethod(currency, accountType, bankCode string) (HaiPayPayoutMethod, bool) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	accountType = strings.ToUpper(strings.TrimSpace(accountType))
	bankCode = strings.TrimSpace(bankCode)
	for _, method := range haiPayPayoutByCurrency[currency].Methods {
		if method.AccountType == accountType && strings.EqualFold(method.BankCode, bankCode) {
			return method, true
		}
	}
	return HaiPayPayoutMethod{}, false
}

func IsHaiPayPayoutMethod(currency, accountType, bankCode string) bool {
	_, ok := FindHaiPayPayoutMethod(currency, accountType, bankCode)
	return ok
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
