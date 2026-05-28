package monthlyloginstatdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var monthlyLoginStatCacheMgr *cache.CacheMgr

func InitMonthlyLoginStatDao() {
	monthlyLoginStatCacheMgr = cache.NewCacheMgr()
}

// GetByMonth 按月标识获取每月登录统计,不存在则新建内存对象
func GetByMonth(month string) *entity.MonthlyLoginStat {
	cacheData := monthlyLoginStatCacheMgr.GetData(month, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.MonthlyLoginStat
		_ = g.Model(string(entity.TbMonthlyLoginStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyLoginStat(month), nil
	})
	return cacheData.(*entity.MonthlyLoginStat)
}
