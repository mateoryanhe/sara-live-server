package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var (
	dailyMiniGamePlusClickStatCacheMgr     *cache.RowCache[*entity.DailyMiniGamePlusClickStat]
)

func initDailyMiniGamePlusClickStatDao() {
	dailyMiniGamePlusClickStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGamePlusClickStat]()
}

func GetDailyMiniGamePlusClickStatByDate(date string) *entity.DailyMiniGamePlusClickStat {
	return dailyMiniGamePlusClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGamePlusClickStat, error) {
		var data *entity.DailyMiniGamePlusClickStat
		_ = g.Model(string(entity.TbDailyMiniGamePlusClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGamePlusClickStat(date), nil
	})
}

func ListRecentDailyMiniGamePlusClickStats(limit int) []*entity.DailyMiniGamePlusClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyMiniGamePlusClickStatsFromDB(limit)
}

func loadRecentDailyMiniGamePlusClickStatsFromDB(limit int) []*entity.DailyMiniGamePlusClickStat {
	list := make([]*entity.DailyMiniGamePlusClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGamePlusClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGamePlusClickStats(list)
	return list
}

func reverseDailyMiniGamePlusClickStats(list []*entity.DailyMiniGamePlusClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGamePlusClickStat(data *entity.DailyMiniGamePlusClickStat) {
	if data == nil || data.ID == "" || dailyMiniGamePlusClickStatCacheMgr == nil {
		return
	}
	dailyMiniGamePlusClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
