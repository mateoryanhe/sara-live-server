package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

func QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	total, data := accountdao.GetUserInfo(req)
	for _, val := range data {
		val.IsAnchor = entity.UserTypeIsAnchor(val.UserType)
		if accountCache := accountdao.GetAccountFromCache(val.OpenId, val.Channel, val.ID); accountCache != nil {
			val.OpenId = accountCache.OpenId
			val.IP = accountCache.IP
			val.RegisterIp = accountCache.RegisterIp
			val.RegisterCountry = accountCache.RegisterCountry
			val.LoginCountry = accountCache.LoginCountry
			val.Channel = accountCache.Channel
			val.PhoneAreaCode = accountCache.PhoneAreaCode
			val.Cancel = accountCache.Cancel
			val.Ban = accountCache.Ban
			val.BanApplyTime = accountCache.BanApplyTime
			val.BanTime = accountCache.BanTime
		}
		if userInfoCache := userinfodao.GetUserInfoFromMemory(val.ID); userInfoCache != nil {
			val.Nickname = userInfoCache.Nickname
			val.Phone = userInfoCache.Phone
			val.Avatar = userInfoCache.Avatar
			val.Remark = userInfoCache.Remark
			val.Gold = userInfoCache.Gold
			val.Diamond = userInfoCache.Diamond
			val.ShareCode = userInfoCache.ShareCode
			val.GuildId = userInfoCache.GuildId
			val.UserType = userInfoCache.UserType
			val.IsAnchor = userInfoCache.IsAnchor()
			val.VipLevel = userInfoCache.VipLevel
			val.LastLoginTime = userInfoCache.LastLoginTime
		}
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
