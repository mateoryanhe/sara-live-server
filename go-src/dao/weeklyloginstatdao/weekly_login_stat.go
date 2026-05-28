package weeklyloginstatdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var weeklyLoginStatCacheMgr *cache.CacheMgr

func InitWeeklyLoginStatDao() {
	weeklyLoginStatCacheMgr = cache.NewCacheMgr()
}

// GetByWeek 按周标识获取每周登录统计,不存在则新建内存对象
func GetByWeek(week string) *entity.WeeklyLoginStat {
	cacheData := weeklyLoginStatCacheMgr.GetData(week, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.WeeklyLoginStat
		_ = g.Model(string(entity.TbWeeklyLoginStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyLoginStat(week), nil
	})
	return cacheData.(*entity.WeeklyLoginStat)
}
