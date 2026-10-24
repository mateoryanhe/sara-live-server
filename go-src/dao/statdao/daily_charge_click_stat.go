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
	dailyChargeClickStatCacheMgr     *cache.RowCache[*entity.DailyChargeClickStat]
)

func initDailyChargeClickStatDao() {
	dailyChargeClickStatCacheMgr = cache.NewRowCache[*entity.DailyChargeClickStat]()
}

func GetDailyChargeClickStatByDate(date string) *entity.DailyChargeClickStat {
	return dailyChargeClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyChargeClickStat, error) {
		var data *entity.DailyChargeClickStat
		_ = g.Model(string(entity.TbDailyChargeClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyChargeClickStat(date), nil
	})
}

func ListRecentDailyChargeClickStats(limit int) []*entity.DailyChargeClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyChargeClickStatsFromDB(limit)
}

func loadRecentDailyChargeClickStatsFromDB(limit int) []*entity.DailyChargeClickStat {
	list := make([]*entity.DailyChargeClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyChargeClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyChargeClickStats(list)
	return list
}

func reverseDailyChargeClickStats(list []*entity.DailyChargeClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyChargeClickStat(data *entity.DailyChargeClickStat) {
	if data == nil || data.ID == "" || dailyChargeClickStatCacheMgr == nil {
		return
	}
	dailyChargeClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
