package userwithdraw

import (
	"strings"

	"xr-game-server/constants/country"
)

// 用户提现代付：部分国家固定为电子钱包，App 不需地址/identifyType 等扩展字段。
var userWithdrawPayoutFixedAccountTypeByCountry = map[string]string{
	"ID": country.HaiPayPayoutAccountTypeEWallet,
	"MY": country.HaiPayPayoutAccountTypeEWallet,
}

func UserWithdrawPayoutFixedAccountType(countryCode string) (string, bool) {
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	fixed, ok := userWithdrawPayoutFixedAccountTypeByCountry[code]
	return fixed, ok
}

func appPayoutExtraRequirementsForUserWithdraw(countryCode, currency, bankCode string) country.HaiPayPayoutExtraRequirements {
	if _, ok := UserWithdrawPayoutFixedAccountType(countryCode); ok {
		return country.HaiPayPayoutExtraRequirements{}
	}
	return country.HaiPayPayoutExtraRequirementsFor(currency, bankCode)
}

func userWithdrawPayoutMethodVisible(countryCode, accountType string) bool {
	fixed, ok := UserWithdrawPayoutFixedAccountType(countryCode)
	if !ok {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(accountType), fixed)
}
