package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity/user"
	"xr-game-server/module/ipgeo"
	"xr-game-server/module/upload"
)

func QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	total, data := accountdao.GetUserInfo(req)
	userIds := make([]uint64, 0, len(data))
	for _, val := range data {
		if val != nil && val.ID > 0 {
			userIds = append(userIds, val.ID)
		}
	}
	locations := userloginlocationdao.GetByUserIds(userIds)
	for _, val := range data {
		val.IsAnchor = entity.UserTypeIsAnchor(val.UserType)
		if accountCache := accountdao.GetAccountFromCache(val.OpenId, val.Channel, val.ID); accountCache != nil {
			val.OpenId = accountCache.OpenId
			val.Channel = accountCache.Channel
			val.PhoneAreaCode = accountCache.PhoneAreaCode
			val.Cancel = accountCache.Cancel
			val.Ban = accountCache.Ban
			val.BanApplyTime = accountCache.BanApplyTime
			val.BanTime = accountCache.BanTime
		}
		if location := locations[val.ID]; location != nil {
			val.IP = location.IP
			val.RegisterIp = location.RegisterIp
			val.RegisterCountry = ipgeo.FormatCountryDisplay(location.RegisterCountry)
			val.LoginCountry = ipgeo.FormatCountryDisplay(location.LoginCountry)
		}
		if userInfoCache := userinfodao.GetUserInfoFromMemory(val.ID); userInfoCache != nil {
			val.Nickname = userInfoCache.Nickname
			val.Phone = userInfoCache.Phone
			val.Avatar = userInfoCache.Avatar
			val.Remark = userInfoCache.Remark
			val.Gold = userInfoCache.Gold
			val.Diamond = userInfoCache.Diamond
			val.ShareCode = userInfoCache.ShareCode
			val.UserType = userInfoCache.UserType
			val.IsAnchor = userInfoCache.IsAnchor()
			val.VipLevel = userInfoCache.VipLevel
			val.LastLoginTime = userInfoCache.LastLoginTime
		}
		val.GuildId = liveroomdao.GetAnchorGuildId(val.ID)
		if userExtCache := userinfodao.GetUserExtFromMemory(val.ID); userExtCache != nil {
			val.CanRank = userExtCache.CanRank
			val.PackageName = userExtCache.PackageName
			val.AppVersion = userExtCache.AppVersion
			val.RechargeWhitelist = userExtCache.RechargeWhitelist
		}
		val.Avatar = upload.ResolveAvatarUrlForUser(val.ID, val.Avatar)
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}
