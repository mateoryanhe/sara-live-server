package userinfodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

var userCumulativeStatCacheMgr *cache.CacheMgr

func initUserCumulativeStatDao() {
	userCumulativeStatCacheMgr = cache.NewCacheMgr()
}

// PreloadUserCumulativeStatToCache 批量预热 user_cumulative_stats 缓存
func PreloadUserCumulativeStatToCache(userIds []uint64) {
	if len(userIds) == 0 || userCumulativeStatCacheMgr == nil {
		return
	}
	ctx := gctx.New()
	rows := make([]*entity.UserCumulativeStat, 0, len(userIds))
	err := g.Model(string(entity.TbUserCumulativeStat)).Ctx(ctx).Unscoped().
		WhereIn(string(db.IdName), userIds).
		Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "preload user cumulative stats failed: %v", err)
		return
	}
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		userCumulativeStatCacheMgr.FlushCache(row.ID, row)
	}
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
