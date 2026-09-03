package userinfo

import (
	"context"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/accountdto"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
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
		Account:        toUserAccountDetailItem(account),
		Profile:        toUserProfileDetailItem(userInfo, req.UserId),
		Wallet:         toUserWalletDetailItem(userInfo),
		UserExt:        toUserExtDetailItem(userinfodao.GetUserExtFromDB(req.UserId)),
		CumulativeStat: toUserCumulativeStatDetailItem(userinfodao.GetUserCumulativeStatFromDB(req.UserId)),
		LoginDevice:    toUserLoginDeviceDetailItem(userlogindevicedao.GetUserLoginDeviceFromDB(req.UserId)),
	}, nil
}

func toUserAccountDetailItem(account *userentity.Account) *accountdto.UserAccountDetailItem {
	if account == nil {
		return nil
	}
	item := &accountdto.UserAccountDetailItem{
		ID:              account.ID,
		OpenId:          account.OpenId,
		IP:              account.IP,
		RegisterIp:      account.RegisterIp,
		RegisterCountry: account.RegisterCountry,
		LoginCountry:    account.LoginCountry,
		Channel:         account.Channel,
		PhoneAreaCode:   account.PhoneAreaCode,
		Ban:             account.Ban,
		BanTime:         account.BanTime,
		BanApplyTime:    account.BanApplyTime,
		Cancel:          account.Cancel,
	}
	if !account.CreatedAt.IsZero() {
		createdAt := account.CreatedAt
		item.CreatedAt = &createdAt
	}
	return item
}

func toUserProfileDetailItem(info *userentity.UserInfo, userId uint64) *accountdto.UserProfileDetailItem {
	if info == nil || info.ID == 0 {
		return &accountdto.UserProfileDetailItem{
			GuildId: liveroomdao.GetAnchorGuildId(userId),
		}
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
	if !info.UpdatedAt.IsZero() {
		updatedAt := info.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
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
	}
	if !ext.UpdatedAt.IsZero() {
		updatedAt := ext.UpdatedAt
		item.UpdatedAt = &updatedAt
	}
	return item
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
