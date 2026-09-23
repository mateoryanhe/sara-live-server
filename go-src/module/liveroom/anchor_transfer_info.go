package liveroom

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/country"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

func validatePlatformAnchor(anchorId uint64) bool {
	if anchorId == 0 {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(anchorId)
	room := liveroomdao.ResolveRoom(anchorId)
	return user != nil && userentity.UserTypeIsAnchor(user.UserType) && room != nil && room.GuildId == 0
}

// GetAnchorTransferInfo 获取平台主播独立的收款信息。
func GetAnchorTransferInfo(_ context.Context, req *accountdto.GetAnchorTransferInfoReq) (*accountdto.GetAnchorTransferInfoRes, error) {
	if req == nil || !validatePlatformAnchor(req.AnchorId) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &accountdto.GetAnchorTransferInfoRes{
		Countries: buildAnchorTransferCountryOptions(),
	}
	row := liveroomdao.GetAnchorTransferInfo(req.AnchorId)
	if row == nil {
		res.Info = &accountdto.AnchorTransferInfoItem{AnchorId: strconv.FormatUint(req.AnchorId, 10)}
		return res, nil
	}
	res.Info = toAnchorTransferInfoItem(row)
	return res, nil
}

// SaveAnchorTransferInfo 保存平台主播独立的收款信息。
func SaveAnchorTransferInfo(_ context.Context, req *accountdto.SaveAnchorTransferInfoReq) (*accountdto.SaveAnchorTransferInfoRes, error) {
	if req == nil || !validatePlatformAnchor(req.AnchorId) {
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

	accountType := strings.ToUpper(strings.TrimSpace(req.AccountType))
	method, ok := country.FindHaiPayPayoutMethod(currency, accountType, req.BankCode)
	if !ok {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if strings.TrimSpace(req.PayeeName) == "" || strings.TrimSpace(req.Phone) == "" ||
		strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.AccountNo) == "" ||
		strings.TrimSpace(req.BankCode) == "" {
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
		address1 = ""
		address2 = ""
		address3 = ""
		postalCode = ""
	}

	row := liveentity.NewLiveAnchorTransferInfo(req.AnchorId)
	if existing := liveroomdao.GetAnchorTransferInfo(req.AnchorId); existing != nil {
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
	if err := liveroomdao.SaveAnchorTransferInfo(row); err != nil {
		return nil, err
	}
	return &accountdto.SaveAnchorTransferInfoRes{Success: true}, nil
}

func toAnchorTransferInfoItem(row *liveentity.LiveAnchorTransferInfo) *accountdto.AnchorTransferInfoItem {
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
	return &accountdto.AnchorTransferInfoItem{
		AnchorId:     strconv.FormatUint(row.ID, 10),
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

func buildAnchorTransferCountryOptions() []*accountdto.AnchorTransferCountryOption {
	options := country.ListHaiPayPayoutOptions()
	version := countryflagdeploy.CurrentVersion()
	out := make([]*accountdto.AnchorTransferCountryOption, 0, len(options))
	for _, option := range options {
		c, ok := country.Get(option.CountryCode)
		if !ok {
			c = country.Country{Code: option.CountryCode, NameEn: option.CountryCode}
		}
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		walletOptions := make([]accountdto.AnchorTransferWalletOption, 0, len(option.Wallets))
		for _, wallet := range option.Wallets {
			walletOptions = append(walletOptions, accountdto.AnchorTransferWalletOption{Code: wallet.Code, Name: wallet.Name})
		}
		methodOptions := make([]accountdto.AnchorTransferMethodOption, 0, len(option.Methods))
		for _, method := range option.Methods {
			extra := country.HaiPayPayoutExtraRequirementsFor(option.Currency, method.BankCode)
			methodOptions = append(methodOptions, accountdto.AnchorTransferMethodOption{
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
		out = append(out, &accountdto.AnchorTransferCountryOption{
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
