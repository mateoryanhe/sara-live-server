package stat

import (
	"xr-game-server/dao/userinfodao"
	userentity "xr-game-server/entity/user"
)

func shouldCountUserStat(userId uint64) bool {
	if userId == 0 {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	return !userentity.UserTypeExcludedFromStat(user.UserType)
}
