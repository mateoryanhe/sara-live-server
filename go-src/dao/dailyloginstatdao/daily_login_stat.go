package dailyloginstatdao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var dailyLoginStatCacheMgr *cache.CacheMgr

func InitDailyLoginStatDao() {
	dailyLoginStatCacheMgr = cache.NewCacheMgr()
}

// GetByDate 按日期获取每日登录统计,不存在则新建内存对象
func GetByDate(date string) *entity.DailyLoginStat {
	cacheData := dailyLoginStatCacheMgr.GetData(date, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.DailyLoginStat
		_ = g.Model(string(entity.TbDailyLoginStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyLoginStat(date), nil
	})
	return cacheData.(*entity.DailyLoginStat)
}
