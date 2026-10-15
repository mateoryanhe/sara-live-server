package userwithdraw

import (
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	rechargeentity "xr-game-server/entity/recharge"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
)

// enabledUserWithdrawPayoutCountryCfg 返回 CMS 已启用且配置完整的提现代付国家。
func enabledUserWithdrawPayoutCountryCfg(countryCode string) (*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg, bool) {
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	if countryCode == "" {
		return nil, false
	}
	row := cfgdao.GetHaiPayUserWithdrawPayoutCountryCfgCached(countryCode)
	if row == nil || !row.Enabled || !cfgdao.HaiPayUserWithdrawPayoutCountryCfgComplete(row) {
		return nil, false
	}
	return row, true
}

func enabledUserWithdrawPayoutCountryCfgByCurrency(currency string) (*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg, bool) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return nil, false
	}
	countryCode := country.HaiPayPayoutCountryFromCurrency(currency)
	if countryCode == "" {
		return nil, false
	}
	cfg, ok := enabledUserWithdrawPayoutCountryCfg(countryCode)
	if !ok {
		return nil, false
	}
	if strings.ToUpper(strings.TrimSpace(cfg.CurrencyCode)) != currency {
		return nil, false
	}
	return cfg, true
}

func validateAppWithdrawCurrency(currency string) error {
	if _, ok := enabledUserWithdrawPayoutCountryCfgByCurrency(currency); !ok {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func validateUserPayoutInfoForWithdrawCurrency(info *userentity.UserPayoutInfo, currency string) error {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	countryCode := country.HaiPayPayoutCountryFromCurrency(currency)
	if countryCode == "" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return validateUserPayoutInfoForWithdrawCountry(info, countryCode)
}

func validateUserPayoutInfoForWithdrawCountry(info *userentity.UserPayoutInfo, countryCode string) error {
	if info == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	if countryCode == "" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	cfg, ok := enabledUserWithdrawPayoutCountryCfg(countryCode)
	if !ok {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	infoCountry := strings.ToUpper(strings.TrimSpace(info.CountryCode))
	if infoCountry == "" {
		infoCountry = country.HaiPayPayoutCountryFromCurrency(strings.ToUpper(strings.TrimSpace(info.Currency)))
	}
	if infoCountry != countryCode {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	currency := strings.ToUpper(strings.TrimSpace(cfg.CurrencyCode))
	if strings.ToUpper(strings.TrimSpace(info.Currency)) != currency {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	accountType := strings.ToUpper(strings.TrimSpace(cfg.AccountType))
	if strings.ToUpper(strings.TrimSpace(info.AccountType)) != accountType {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if !strings.EqualFold(strings.TrimSpace(info.BankCode), strings.TrimSpace(cfg.BankCode)) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if _, ok := country.FindHaiPayPayoutMethod(currency, accountType, cfg.BankCode); !ok {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func payoutMethodFromUserWithdrawCountryCfg(countryCode string) (currencyCode, accountType, bankCode string, err error) {
	cfg, ok := enabledUserWithdrawPayoutCountryCfg(countryCode)
	if !ok {
		return "", "", "", errercode.CreateCode(errercode.InvalidParam)
	}
	currencyCode = strings.ToUpper(strings.TrimSpace(cfg.CurrencyCode))
	accountType = strings.ToUpper(strings.TrimSpace(cfg.AccountType))
	bankCode = strings.TrimSpace(cfg.BankCode)
	if _, ok := country.FindHaiPayPayoutMethod(currencyCode, accountType, bankCode); !ok {
		return "", "", "", errercode.CreateCode(errercode.InvalidParam)
	}
	return currencyCode, accountType, bankCode, nil
}
