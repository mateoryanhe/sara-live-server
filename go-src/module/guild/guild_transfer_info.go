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
			GuildId:     strconv.FormatUint(req.GuildId, 10),
			AccountType: liveentity.GuildTransferAccountTypeBank,
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
	currency := country.HaiPayPayoutCurrency(countryCode)
	if currency == "" || !country.IsHaiPayPayoutCountry(countryCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildByIdFromDB(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	accountType := strings.TrimSpace(req.AccountType)
	if accountType == "" {
		accountType = liveentity.GuildTransferAccountTypeBank
	}

	row := liveentity.NewLiveGuildTransferInfo(req.GuildId)
	if existing := guilddao.GetGuildTransferInfo(req.GuildId); existing != nil {
		row.CreatedAt = existing.CreatedAt
	}
	row.CountryCode = countryCode
	row.Currency = currency
	row.AccountType = accountType
	row.PayeeName = strings.TrimSpace(req.PayeeName)
	row.Phone = strings.TrimSpace(req.Phone)
	row.Email = strings.ToLower(strings.TrimSpace(req.Email))
	row.BankName = strings.TrimSpace(req.BankName)
	row.AccountNo = strings.TrimSpace(req.AccountNo)
	row.BankCode = strings.TrimSpace(req.BankCode)
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
	accountType := strings.TrimSpace(row.AccountType)
	if accountType == "" {
		accountType = liveentity.GuildTransferAccountTypeBank
	}
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
	regions := country.ListHaiPayPayoutCountries()
	version := countryflagdeploy.CurrentVersion()
	out := make([]*guilddto.GuildTransferCountryOption, 0, len(regions))
	for _, c := range regions {
		cur := country.HaiPayPayoutCurrency(c.Code)
		if cur == "" {
			continue
		}
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		out = append(out, &guilddto.GuildTransferCountryOption{
			CountryCode: c.Code,
			NameEn:      c.NameEn,
			NameZh:      c.NameZh,
			Currency:    cur,
			Icon:        icon,
		})
	}
	return out
}
