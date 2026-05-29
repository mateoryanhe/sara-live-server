package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
)

func QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	total, data := accountdao.GetUserInfo(req)
	//检查缓存数据问题
	//保证cms系统查询数据正确
	for _, val := range data {
		if InCache(val.ID) {
			//TODO:记录日志
			userInfoCache := userinfodao.GetUserInfoByUserId(val.ID)
			val.Diamond = userInfoCache.Diamond
			val.Gold = userInfoCache.Gold
			val.IsAnchor = userInfoCache.IsAnchor
			accountCache := accountdao.GetAccountBy(val.OpenId, val.Channel)
			val.Cancel = accountCache.Cancel
			val.Ban = accountCache.Ban
			val.BanApplyTime = accountCache.BanApplyTime
			val.BanTime = accountCache.BanTime
			val.VipLevel = userInfoCache.VipLevel
		}
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}
