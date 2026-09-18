package guild

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/country"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

func GetGuildTransferInfo(_ context.Context, req *guilddto.GetGuildTransferInfoReq) (*guilddto.GetGuildTransferInfoRes, error) {
	if req == nil || req.GuildId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildByIdFromDB(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	res := &guilddto.GetGuildTransferInfoRes{
		Countries: buildGuildTransferCountryOptions(),
	}
	row := guilddao.GetGuildTransferInfo(req.GuildId)
	if row == nil {
		res.Info = &guilddto.GuildTransferInfoItem{
			GuildId: strconv.FormatUint(req.GuildId, 10),
		}
		return res, nil
	}
	res.Info = toGuildTransferInfoItem(row)
	return res, nil
}

func SaveGuildTransferInfo(_ context.Context, req *guilddto.SaveGuildTransferInfoReq) (*guilddto.SaveGuildTransferInfoRes, error) {
	if req == nil || req.GuildId == 0 {
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
	guild := guilddao.GetGuildByIdFromDB(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
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

	row := liveentity.NewLiveGuildTransferInfo(req.GuildId)
	if existing := guilddao.GetGuildTransferInfo(req.GuildId); existing != nil {
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
	row.Remark = strings.TrimSpace(req.Remark)
	row.UpdatedAt = time.Now()

	if err := guilddao.SaveGuildTransferInfo(row); err != nil {
		return nil, err
	}
	return &guilddto.SaveGuildTransferInfoRes{Success: true}, nil
}

func toGuildTransferInfoItem(row *liveentity.LiveGuildTransferInfo) *guilddto.GuildTransferInfoItem {
	if row == nil {
		return nil
	}
	updatedAt := ""
	if !row.UpdatedAt.IsZero() {
		updatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	countryCode := strings.ToUpper(strings.TrimSpace(row.CountryCode))
	currency := strings.ToUpper(strings.TrimSpace(row.Currency))
	// 兼容旧数据:仅有币种时反推国家
	if countryCode == "" && currency != "" {
		countryCode = country.HaiPayPayoutCountryFromCurrency(currency)
	}
	if currency == "" && countryCode != "" {
		currency = country.HaiPayPayoutCurrency(countryCode)
	}
	accountType := strings.ToUpper(strings.TrimSpace(row.AccountType))
	return &guilddto.GuildTransferInfoItem{
		GuildId:     strconv.FormatUint(row.ID, 10),
		CountryCode: countryCode,
		Currency:    currency,
		AccountType: accountType,
		PayeeName:   row.PayeeName,
		Phone:       row.Phone,
		Email:       row.Email,
		BankName:    row.BankName,
		AccountNo:   row.AccountNo,
		BankCode:    row.BankCode,
		Remark:      row.Remark,
		UpdatedAt:   updatedAt,
	}
}

func buildGuildTransferCountryOptions() []*guilddto.GuildTransferCountryOption {
	options := country.ListHaiPayPayoutOptions()
	version := countryflagdeploy.CurrentVersion()
	out := make([]*guilddto.GuildTransferCountryOption, 0, len(options))
	for _, option := range options {
		c, ok := country.Get(option.CountryCode)
		if !ok {
			c = country.Country{Code: option.CountryCode, NameEn: option.CountryCode}
		}
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		walletOptions := make([]guilddto.GuildTransferWalletOption, 0, len(option.Wallets))
		for _, wallet := range option.Wallets {
			walletOptions = append(walletOptions, guilddto.GuildTransferWalletOption{
				Code: wallet.Code,
				Name: wallet.Name,
			})
		}
		methodOptions := make([]guilddto.GuildTransferMethodOption, 0, len(option.Methods))
		for _, method := range option.Methods {
			methodOptions = append(methodOptions, guilddto.GuildTransferMethodOption{
				AccountType: method.AccountType,
				BankCode:    method.BankCode,
				Limit:       method.Limit,
				Description: method.Description,
			})
		}
		out = append(out, &guilddto.GuildTransferCountryOption{
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
