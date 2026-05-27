package userinfo

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/module/upload"
)

// QueryAnchorList CMS分页查询主播列表
func QueryAnchorList(_ context.Context, req *accountdto.QueryAnchorListReq) (*httpserver.CMSQueryResp, error) {
	total, data := accountdao.GetAnchorList(req)
	for _, val := range data {
		if val == nil {
			continue
		}
		val.Avatar = upload.GetUrlByName(val.Avatar)
		if InCache(val.ID) {
			userInfoCache := userinfodao.GetUserInfoByUserId(val.ID)
			if userInfoCache != nil {
				val.Nickname = userInfoCache.Nickname
				val.Phone = userInfoCache.Phone
				val.GuildId = userInfoCache.GuildId
			}
			accountCache := accountdao.GetAccountById(val.ID)
			if accountCache != nil {
				val.IP = accountCache.IP
			}
		}
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}
