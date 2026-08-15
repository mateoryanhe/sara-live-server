package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyLoginStatCacheMgr *cache.CacheMgr

func initMonthlyLoginStatDao() {
	monthlyLoginStatCacheMgr = cache.NewCacheMgr()
}

// GetMonthlyLoginStatByMonth 按月标识获取每月登录统计,不存在则新建内存对象
func GetMonthlyLoginStatByMonth(month string) *entity.MonthlyLoginStat {
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

// ListRecentMonthlyLoginStats 查询最近N月登录统计(按时间正序)
func ListRecentMonthlyLoginStats(limit int) []*entity.MonthlyLoginStat {
	list := make([]*entity.MonthlyLoginStat, 0)
	if limit <= 0 {
		limit = 12
	}
	_ = g.Model(string(entity.TbMonthlyLoginStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyLoginStats(list)
	return list
}

func reverseMonthlyLoginStats(list []*entity.MonthlyLoginStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}
