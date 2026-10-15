package userwithdraw

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/country"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/userwithdrawdao"
	"xr-game-server/dto/userwithdrawdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

func GetPayoutInfo(_ context.Context, userId uint64) (*userwithdrawdto.AppGetPayoutInfoRes, error) {
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &userwithdrawdto.AppGetPayoutInfoRes{
		Countries: buildPayoutCountryOptions(),
		Info:      &userwithdrawdto.AppPayoutInfoItem{},
	}
	row := userwithdrawdao.GetUserPayoutInfo(userId)
	if row != nil {
		res.Info = toPayoutInfoItem(row)
	}
	return res, nil
}

func SavePayoutInfo(_ context.Context, userId uint64, req *userwithdrawdto.AppSavePayoutInfoReq) (*userwithdrawdto.AppSavePayoutInfoRes, error) {
	if userId == 0 || req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	countryCode := strings.ToUpper(strings.TrimSpace(req.CountryCode))
	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency != "" {
		currencyCountry := country.HaiPayPayoutCountryFromCurrency(currency)
		if currencyCountry == "" || (countryCode != "" && countryCode != currencyCountry) {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		countryCode = currencyCountry
	} else {
		currency = country.HaiPayPayoutCurrency(countryCode)
	}
	if currency == "" || !country.IsHaiPayPayoutCountry(countryCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	cfg, ok := enabledUserWithdrawPayoutCountryCfg(countryCode)
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	currency = strings.ToUpper(strings.TrimSpace(cfg.CurrencyCode))
	accountType := strings.ToUpper(strings.TrimSpace(cfg.AccountType))
	bankCode := strings.TrimSpace(cfg.BankCode)
	method, ok := country.FindHaiPayPayoutMethod(currency, accountType, bankCode)
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if strings.TrimSpace(req.PayeeName) == "" || strings.TrimSpace(req.Phone) == "" ||
		strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.AccountNo) == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if reqAccountType := strings.ToUpper(strings.TrimSpace(req.AccountType)); reqAccountType != "" && reqAccountType != accountType {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if reqBank := strings.TrimSpace(req.BankCode); reqBank != "" && !strings.EqualFold(reqBank, bankCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	extraRequirements := country.HaiPayPayoutExtraRequirementsFor(currency, method.BankCode)
	identifyType := strings.TrimSpace(req.IdentifyType)
	if len(extraRequirements.IdentifyTypeOptions) > 0 {
		identifyType = strings.ToUpper(identifyType)
	}
	if !extraRequirements.IsAllowedIdentifyType(identifyType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	address1 := strings.TrimSpace(req.Address1)
	address2 := strings.TrimSpace(req.Address2)
	address3 := strings.TrimSpace(req.Address3)
	postalCode := strings.TrimSpace(req.PostalCode)
	if extraRequirements.AddressRequired &&
		(address1 == "" || address2 == "" || address3 == "" || postalCode == "") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !extraRequirements.IdentifyTypeRequired {
		identifyType = ""
	}
	if !extraRequirements.AddressRequired {
		address1, address2, address3, postalCode = "", "", "", ""
	}

	row := userentity.NewUserPayoutInfo(userId)
	if existing := userwithdrawdao.GetUserPayoutInfo(userId); existing != nil {
		row.CreatedAt = existing.CreatedAt
	}
	row.CountryCode = countryCode
	row.Currency = currency
	row.AccountType = method.AccountType
	row.PayeeName = strings.TrimSpace(req.PayeeName)
	row.Phone = strings.TrimSpace(req.Phone)
	row.Email = strings.ToLower(strings.TrimSpace(req.Email))
	row.BankName = method.Description
	row.AccountNo = strings.TrimSpace(req.AccountNo)
	row.BankCode = method.BankCode
	row.IdentifyType = identifyType
	row.Address1 = address1
	row.Address2 = address2
	row.Address3 = address3
	row.PostalCode = postalCode
	row.Remark = strings.TrimSpace(req.Remark)
	row.UpdatedAt = time.Now()
	if err := userwithdrawdao.SaveUserPayoutInfo(row); err != nil {
		return nil, err
	}
	return &userwithdrawdto.AppSavePayoutInfoRes{Success: true}, nil
}

func toPayoutInfoItem(row *userentity.UserPayoutInfo) *userwithdrawdto.AppPayoutInfoItem {
	if row == nil {
		return nil
	}
	updatedAt := ""
	if !row.UpdatedAt.IsZero() {
		updatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	countryCode := strings.ToUpper(strings.TrimSpace(row.CountryCode))
	currency := strings.ToUpper(strings.TrimSpace(row.Currency))
	if countryCode == "" && currency != "" {
		countryCode = country.HaiPayPayoutCountryFromCurrency(currency)
	}
	if currency == "" && countryCode != "" {
		currency = country.HaiPayPayoutCurrency(countryCode)
	}
	return &userwithdrawdto.AppPayoutInfoItem{
		CountryCode:  countryCode,
		Currency:     currency,
		AccountType:  strings.ToUpper(strings.TrimSpace(row.AccountType)),
		PayeeName:    row.PayeeName,
		Phone:        row.Phone,
		Email:        row.Email,
		BankName:     row.BankName,
		AccountNo:    row.AccountNo,
		BankCode:     row.BankCode,
		IdentifyType: row.IdentifyType,
		Address1:     row.Address1,
		Address2:     row.Address2,
		Address3:     row.Address3,
		PostalCode:   row.PostalCode,
		Remark:       row.Remark,
		UpdatedAt:    updatedAt,
	}
}

func buildPayoutCountryOptions() []*userwithdrawdto.AppPayoutCountryOption {
	configured := cfgdao.GetHaiPayUserWithdrawPayoutCountryCfgMapCached()
	options := listAppVisibleHaiPayPayoutOptions()
	version := countryflagdeploy.CurrentVersion()
	out := make([]*userwithdrawdto.AppPayoutCountryOption, 0, len(options))
	for _, option := range options {
		c, ok := country.Get(option.CountryCode)
		if !ok {
			c = country.Country{Code: option.CountryCode, NameEn: option.CountryCode}
		}
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		walletOptions := make([]userwithdrawdto.AppPayoutWalletOption, 0, len(option.Wallets))
		for _, wallet := range option.Wallets {
			walletOptions = append(walletOptions, userwithdrawdto.AppPayoutWalletOption{Code: wallet.Code, Name: wallet.Name})
		}
		methodOptions := make([]userwithdrawdto.AppPayoutMethodOption, 0, len(option.Methods))
		stored := configured[option.CountryCode]
		for _, method := range option.Methods {
			if stored != nil && cfgdao.HaiPayUserWithdrawPayoutCountryCfgComplete(stored) {
				if method.AccountType != strings.ToUpper(strings.TrimSpace(stored.AccountType)) ||
					!strings.EqualFold(method.BankCode, stored.BankCode) {
					continue
				}
			}
			extra := country.HaiPayPayoutExtraRequirementsFor(option.Currency, method.BankCode)
			methodOptions = append(methodOptions, userwithdrawdto.AppPayoutMethodOption{
				AccountType:          method.AccountType,
				BankCode:             method.BankCode,
				Limit:                method.Limit,
				Description:          method.Description,
				IdentifyTypeRequired: extra.IdentifyTypeRequired,
				IdentifyTypeOptions:  append([]string(nil), extra.IdentifyTypeOptions...),
				CountryRequired:      extra.CountryRequired,
				AddressRequired:      extra.AddressRequired,
			})
		}
		out = append(out, &userwithdrawdto.AppPayoutCountryOption{
			CountryCode:  c.Code,
			NameEn:       c.NameEn,
			NameZh:       c.NameZh,
			Currency:     option.Currency,
			Region:       option.Region,
			AccountTypes: append([]string(nil), option.AccountTypes...),
			Wallets:      walletOptions,
			Methods:      methodOptions,
			Icon:         icon,
		})
	}
	return out
}

func formatLogTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatLogTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func parseIdList(ids []string) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, raw := range ids {
		id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
		if err != nil || id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
