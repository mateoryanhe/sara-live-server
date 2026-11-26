package userattributionreportdao

import (
	"github.com/gogf/gf/v2/frame/g"
	userentity "xr-game-server/entity/user"
)

const cmsListLimit = 100

// ListByUserId CMS 用户详情：按 id 倒序查 AF 归因上报(append-only)
func ListByUserId(userId uint64) []*userentity.UserAttributionReport {
	if userId == 0 {
		return nil
	}
	var list []*userentity.UserAttributionReport
	err := g.DB().Model(string(userentity.TbUserAttributionReport)).
		Where("user_id", userId).
		Order("id desc").
		Limit(cmsListLimit).
		Scan(&list)
	if err != nil {
		return nil
	}
	return list
}