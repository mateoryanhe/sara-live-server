package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyLoginStatCacheMgr *cache.CacheMgr

func initWeeklyLoginStatDao() {
	weeklyLoginStatCacheMgr = cache.NewCacheMgr()
}

// GetWeeklyLoginStatByWeek 按周标识获取每周登录统计,不存在则新建内存对象
func GetWeeklyLoginStatByWeek(week string) *entity.WeeklyLoginStat {
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

// ListRecentWeeklyLoginStats 查询最近N周登录统计(按时间正序)
func ListRecentWeeklyLoginStats(limit int) []*entity.WeeklyLoginStat {
	list := make([]*entity.WeeklyLoginStat, 0)
	if limit <= 0 {
		limit = 12
	}
	_ = g.Model(string(entity.TbWeeklyLoginStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyLoginStats(list)
	return list
}

func reverseWeeklyLoginStats(list []*entity.WeeklyLoginStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}
