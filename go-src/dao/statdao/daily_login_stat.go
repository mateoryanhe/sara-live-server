package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyLoginStatCacheMgr *cache.RowCache[*entity.DailyLoginStat]

func initDailyLoginStatDao() {
	dailyLoginStatCacheMgr = cache.NewRowCache[*entity.DailyLoginStat]()
}

// GetDailyLoginStatByDate 按日期获取每日登录统计,不存在则新建内存对象
func GetDailyLoginStatByDate(date string) *entity.DailyLoginStat {
	return dailyLoginStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyLoginStat, error) {
		var data *entity.DailyLoginStat
		_ = g.Model(string(entity.TbDailyLoginStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyLoginStat(date), nil
	})
}

// ListRecentDailyLoginStats 查询最近N天登录统计(按时间正序)
func ListRecentDailyLoginStats(limit int) []*entity.DailyLoginStat {
	list := make([]*entity.DailyLoginStat, 0)
	if limit <= 0 {
		limit = 30
	}
	_ = g.Model(string(entity.TbDailyLoginStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyLoginStats(list)
	return list
}

func reverseDailyLoginStats(list []*entity.DailyLoginStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}
