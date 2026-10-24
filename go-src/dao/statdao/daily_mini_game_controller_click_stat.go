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
	dailyMiniGameControllerClickStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameControllerClickStat]
)

func initDailyMiniGameControllerClickStatDao() {
	dailyMiniGameControllerClickStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameControllerClickStat]()
}

func GetDailyMiniGameControllerClickStatByDate(date string) *entity.DailyMiniGameControllerClickStat {
	return dailyMiniGameControllerClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameControllerClickStat, error) {
		var data *entity.DailyMiniGameControllerClickStat
		_ = g.Model(string(entity.TbDailyMiniGameControllerClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameControllerClickStat(date), nil
	})
}

func ListRecentDailyMiniGameControllerClickStats(limit int) []*entity.DailyMiniGameControllerClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyMiniGameControllerClickStatsFromDB(limit)
}

func loadRecentDailyMiniGameControllerClickStatsFromDB(limit int) []*entity.DailyMiniGameControllerClickStat {
	list := make([]*entity.DailyMiniGameControllerClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameControllerClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameControllerClickStats(list)
	return list
}

func reverseDailyMiniGameControllerClickStats(list []*entity.DailyMiniGameControllerClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameControllerClickStat(data *entity.DailyMiniGameControllerClickStat) {
	if data == nil || data.ID == "" || dailyMiniGameControllerClickStatCacheMgr == nil {
		return
	}
	dailyMiniGameControllerClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
