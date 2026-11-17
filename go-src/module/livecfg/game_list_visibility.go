package livecfg

import (
	"xr-game-server/dao/userinfodao"
	userentity "xr-game-server/entity/user"
)

// ShouldHideAppGameListForUser 按直播配置与用户身份判断是否返回空游戏列表。
func ShouldHideAppGameListForUser(userId uint64) bool {
	info := userinfodao.GetUserInfoByUserId(userId)
	isAnchor := info != nil && userentity.UserTypeIsAnchor(info.UserType)
	if isAnchor {
		return IsAnchorGameListHidden()
	}
	return IsAudienceGameListHidden()
}