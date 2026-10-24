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
	dailyMiniGameSpinClickStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameSpinClickStat]
)

func initDailyMiniGameSpinClickStatDao() {
	dailyMiniGameSpinClickStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameSpinClickStat]()
}

func GetDailyMiniGameSpinClickStatByDate(date string) *entity.DailyMiniGameSpinClickStat {
	return dailyMiniGameSpinClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameSpinClickStat, error) {
		var data *entity.DailyMiniGameSpinClickStat
		_ = g.Model(string(entity.TbDailyMiniGameSpinClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameSpinClickStat(date), nil
	})
}

func ListRecentDailyMiniGameSpinClickStats(limit int) []*entity.DailyMiniGameSpinClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyMiniGameSpinClickStatsFromDB(limit)
}

func loadRecentDailyMiniGameSpinClickStatsFromDB(limit int) []*entity.DailyMiniGameSpinClickStat {
	list := make([]*entity.DailyMiniGameSpinClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameSpinClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameSpinClickStats(list)
	return list
}

func reverseDailyMiniGameSpinClickStats(list []*entity.DailyMiniGameSpinClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameSpinClickStat(data *entity.DailyMiniGameSpinClickStat) {
	if data == nil || data.ID == "" || dailyMiniGameSpinClickStatCacheMgr == nil {
		return
	}
	dailyMiniGameSpinClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
