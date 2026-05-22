package userinfodao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var userCumulativeStatCacheMgr *cache.CacheMgr

func initUserCumulativeStatDao() {
	userCumulativeStatCacheMgr = cache.NewCacheMgr()
}

// GetUserCumulativeStatByUserId 根据玩家ID获取累计数值,命中不了缓存从数据库拉取,数据库不存在则新建
func GetUserCumulativeStatByUserId(userId uint64) *entity.UserCumulativeStat {
	cacheData := userCumulativeStatCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.UserCumulativeStat
		err = g.Model(string(entity.TbUserCumulativeStat)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewUserCumulativeStat(userId), nil
	})
	return cacheData.(*entity.UserCumulativeStat)
}
