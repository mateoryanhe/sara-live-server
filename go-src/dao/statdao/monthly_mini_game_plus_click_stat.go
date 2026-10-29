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
	monthlyMiniGamePlusClickStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGamePlusClickStat]
)

func initMonthlyMiniGamePlusClickStatDao() {
	monthlyMiniGamePlusClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGamePlusClickStat]()
}

func GetMonthlyMiniGamePlusClickStatByMonth(month string) *entity.MonthlyMiniGamePlusClickStat {
	return monthlyMiniGamePlusClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGamePlusClickStat, error) {
		var data *entity.MonthlyMiniGamePlusClickStat
		_ = g.Model(string(entity.TbMonthlyMiniGamePlusClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGamePlusClickStat(month), nil
	})
}

func ListRecentMonthlyMiniGamePlusClickStats(limit int) []*entity.MonthlyMiniGamePlusClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyMiniGamePlusClickStatsFromDB(limit)
}

func loadRecentMonthlyMiniGamePlusClickStatsFromDB(limit int) []*entity.MonthlyMiniGamePlusClickStat {
	list := make([]*entity.MonthlyMiniGamePlusClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGamePlusClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGamePlusClickStats(list)
	return list
}

func reverseMonthlyMiniGamePlusClickStats(list []*entity.MonthlyMiniGamePlusClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGamePlusClickStat(data *entity.MonthlyMiniGamePlusClickStat) {
	if data == nil || data.ID == "" || monthlyMiniGamePlusClickStatCacheMgr == nil {
		return
	}
	monthlyMiniGamePlusClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
