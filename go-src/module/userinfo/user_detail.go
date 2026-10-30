package userinfo

import (
	"context"
	"strings"

	"xr-game-server/constants/country"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dao/channelpaydao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/dao/userwithdrawdao"
	"xr-game-server/dto/accountdto"
	rechargeentity "xr-game-server/entity/recharge"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/ipgeo"
	"xr-game-server/module/upload"
)

// QueryUserDetail CMS获取用户详情
func QueryUserDetail(_ context.Context, req *accountdto.GetUserDetailReq) (*accountdto.GetUserDetailRes, error) {
	account := accountdao.GetAccountById(req.UserId)
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userInfo := userinfodao.GetUserInfoFromDB(req.UserId)
	return &accountdto.GetUserDetailRes{
		Account:        toUserAccountDetailItem(account, userloginlocationdao.GetByUserId(req.UserId)),
		Profile:        toUserProfileDetailItem(userInfo, req.UserId),
		Wallet:         toUserWalletDetailItem(userInfo),
		UserExt:        toUserExtDetailItem(userinfodao.GetUserExtFromDB(req.UserId)),
		CumulativeStat: toUserCumulativeStatDetailItem(userinfodao.GetUserCumulativeStatFromDB(req.UserId)),
		LoginDevice:       toUserLoginDeviceDetailItem(userlogindevicedao.GetUserLoginDeviceFromDB(req.UserId)),
		ChannelPayProfile: toUserChannelPayProfileDetailItem(channelpaydao.GetExisting(req.UserId)),
		WithdrawPayout:    toUserWithdrawPayoutDetailItem(userwithdrawdao.GetUserPayoutInfo(req.UserId)),
	}, nil
}

func toUserWithdrawPayoutDetailItem(row *userentity.UserPayoutInfo) *accountdto.UserWithdrawPayoutDetailItem {
	if row == nil || row.ID == 0 {
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
	return &accountdto.UserWithdrawPayoutDetailItem{
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

func toUserChannelPayProfileDetailItem(row *rechargeentity.ChannelPayUserProfile) *accountdto.UserChannelPayProfileDetailItem {
	if row == nil || row.ID == 0 {
		return &accountdto.UserChannelPayProfileDetailItem{}
	}
	return &accountdto.UserChannelPayProfileDetailItem{
		Name:  strings.TrimSpace(row.Name),
		Email: strings.TrimSpace(row.Email),
	}
}

func toUserAccountDetailItem(account *userentity.Account, location *userentity.UserLoginLocation) *accountdto.UserAccountDetailItem {
	if account == nil {
		return nil
	}
	item := &accountdto.UserAccountDetailItem{
		ID:            account.ID,
		OpenId:        account.OpenId,
		Channel:       account.Channel,
		PhoneAreaCode: account.PhoneAreaCode,
		Ban:           account.Ban,
		BanTime:       account.BanTime,
		BanApplyTime:  account.BanApplyTime,
		Cancel:        account.Cancel,
	}
	if location != nil {
		item.IP = location.IP
		item.RegisterIp = location.RegisterIp
		item.RegisterCountry = ipgeo.FormatCountryDisplay(location.RegisterCountry)
		item.LoginCountry = ipgeo.FormatCountryDisplay(location.LoginCountry)
	}
	if !account.CreatedAt.IsZero() {
		createdAt := account.CreatedAt
		item.CreatedAt = &createdAt
	}
	return item
}

func toUserProfileDetailItem(info *userentity.UserInfo, userId uint64) *accountdto.UserProfileDetailItem {
	if info == nil || info.ID == 0 {
		item := &accountdto.UserProfileDetailItem{
			GuildId: liveroomdao.GetAnchorGuildId(userId),
		}
		fillUserProfileGuildName(item)
		return item
	}
	item := &accountdto.UserProfileDetailItem{
		Nickname:        info.Nickname,
		Phone:           info.Phone,
		Avatar:          upload.ResolveAvatarUrlForUser(info.ID, info.Avatar),
		Remark:          info.Remark,
		ShareCode:       info.ShareCode,
		UserType:        info.UserType,
		IsAnchor:        info.IsAnchor(),
		InviterId:       info.InviterId,
		VipLevel:        info.VipLevel,
		LastLoginTime:   info.LastLoginTime,
		LiveRoomId:      info.LiveRoomId,
		LiveRoomVer:     info.LiveRoomVer,
		Gender:          info.Gender,
		Birthday:        info.Birthday,
		BotAnchorStatus: info.BotAnchorStatus,
		GuildId:         liveroomdao.GetAnchorGuildId(info.ID),
	}
	fillUserProfileGuildName(item)
	if !info.UpdatedAt.IsZero() {
		updatedAt := info.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func fillUserProfileGuildName(item *accountdto.UserProfileDetailItem) {
	if item == nil || item.GuildId == 0 {
		return
	}
	guild := guilddao.GetGuildByIdFromDB(item.GuildId)
	if guild == nil {
		return
	}
	item.GuildName = guild.Name
}

func toUserWalletDetailItem(info *userentity.UserInfo) *accountdto.UserWalletDetailItem {
	if info == nil || info.ID == 0 {
		return &accountdto.UserWalletDetailItem{}
	}
	return &accountdto.UserWalletDetailItem{
		Gold:    info.Gold,
		Diamond: info.Diamond,
	}
}

func toUserExtDetailItem(ext *userentity.UserExt) *accountdto.UserExtDetailItem {
	if ext == nil || ext.ID == 0 {
		return nil
	}
	item := &accountdto.UserExtDetailItem{
		CanRank:                   ext.CanRank,
		PrettyId:                  ext.PrettyId,
		PackageName:               ext.PackageName,
		AppVersion:                ext.AppVersion,
		FollowCount:               ext.FollowCount,
		FollowerCount:             ext.FollowerCount,
		CancelCode:                ext.CancelCode,
		CancelCodeExpireAt:        ext.CancelCodeExpireAt,
		RechargeWhitelist:         ext.RechargeWhitelist,
		ShortVideoUnsettledIncome: ext.ShortVideoUnsettledIncome,
		Email:                     ext.Email,
		ApkPromotionChannel:       strings.TrimSpace(ext.ApkPromotionChannel),
		ApkPromotionChannelName:   resolveApkPromotionChannelName(ext.ApkPromotionChannel),
	}
	if !ext.UpdatedAt.IsZero() {
		updatedAt := ext.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func resolveApkPromotionChannelName(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	row := apkpromotionchanneldao.GetByChannelCode(code)
	if row == nil {
		return ""
	}
	return strings.TrimSpace(row.ChannelName)
}

func toUserCumulativeStatDetailItem(stat *userentity.UserCumulativeStat) *accountdto.UserCumulativeStatDetailItem {
	if stat == nil || stat.ID == 0 {
		return &accountdto.UserCumulativeStatDetailItem{}
	}
	item := &accountdto.UserCumulativeStatDetailItem{
		TotalRecharge:       stat.TotalRechargeGold,
		TotalWithdraw:       stat.TotalWithdraw,
		TotalPayCount:       stat.TotalPayCount,
		TotalDiamondConsume: stat.TotalDiamondConsume,
		TotalGoldConsume:    stat.TotalGoldConsume,
	}
	if !stat.UpdatedAt.IsZero() {
		updatedAt := stat.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}

func toUserLoginDeviceDetailItem(device *userentity.UserLoginDevice) *accountdto.UserLoginDeviceDetailItem {
	if device == nil || device.ID == 0 {
		return nil
	}
	item := &accountdto.UserLoginDeviceDetailItem{
		DeviceType:  device.DeviceType,
		DeviceModel: device.DeviceModel,
		CpuModel:    device.CpuModel,
		OsVersion:   device.OsVersion,
		AppVersion:  device.AppVersion,
		DeviceId:    device.DeviceId,
	}
	if !device.UpdatedAt.IsZero() {
		updatedAt := device.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
}
