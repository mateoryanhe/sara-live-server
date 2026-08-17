package userinfodao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	userentity "xr-game-server/entity/user"
)

// GetUserInfoFromDB 直查 user_infos(不存在返回 nil,不新建)
func GetUserInfoFromDB(userId uint64) *userentity.UserInfo {
	if userId == 0 {
		return nil
	}
	var row userentity.UserInfo
	err := g.Model(string(userentity.TbUserInfo)).Unscoped().WherePri(userId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetUserExtFromDB 直查 user_exts(不存在返回 nil,不新建)
func GetUserExtFromDB(userId uint64) *userentity.UserExt {
	if userId == 0 {
		return nil
	}
	var row userentity.UserExt
	err := g.Model(string(userentity.TbUserExt)).Unscoped().WherePri(userId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetUserCumulativeStatFromDB 直查 user_cumulative_stats(不存在返回 nil,不新建)
func GetUserCumulativeStatFromDB(userId uint64) *userentity.UserCumulativeStat {
	if userId == 0 {
		return nil
	}
	var row userentity.UserCumulativeStat
	err := g.Model(string(userentity.TbUserCumulativeStat)).Unscoped().Where(g.Map{
		string(db.IdName): userId,
	}).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}
