package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userextdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

func QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	total, data := accountdao.GetUserInfo(req)
	//检查缓存数据问题
	//保证cms系统查询数据正确
	for _, val := range data {
		val.IsAnchor = entity.UserTypeIsAnchor(val.UserType)
		if InCache(val.ID) {
			//TODO:记录日志
			userInfoCache := userinfodao.GetUserInfoByUserId(val.ID)
			val.Diamond = userInfoCache.Diamond
			val.Gold = userInfoCache.Gold
			val.IsAnchor = userInfoCache.IsAnchor()
			val.UserType = userInfoCache.UserType
			val.Avatar = userInfoCache.Avatar
			accountCache := accountdao.GetAccountBy(val.OpenId, val.Channel)
			val.Cancel = accountCache.Cancel
			val.Ban = accountCache.Ban
			val.BanApplyTime = accountCache.BanApplyTime
			val.BanTime = accountCache.BanTime
			val.VipLevel = userInfoCache.VipLevel
			val.CanRank = userextdao.GetByUserId(val.ID).CanRank
		}
		val.Avatar = upload.ResolveAvatarUrlForUser(val.ID, val.Avatar)
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}
