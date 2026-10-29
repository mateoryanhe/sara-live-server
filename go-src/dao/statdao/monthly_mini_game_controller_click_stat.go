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
	monthlyMiniGameControllerClickStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameControllerClickStat]
)

func initMonthlyMiniGameControllerClickStatDao() {
	monthlyMiniGameControllerClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameControllerClickStat]()
}

func GetMonthlyMiniGameControllerClickStatByMonth(month string) *entity.MonthlyMiniGameControllerClickStat {
	return monthlyMiniGameControllerClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameControllerClickStat, error) {
		var data *entity.MonthlyMiniGameControllerClickStat
		_ = g.Model(string(entity.TbMonthlyMiniGameControllerClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameControllerClickStat(month), nil
	})
}

func ListRecentMonthlyMiniGameControllerClickStats(limit int) []*entity.MonthlyMiniGameControllerClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyMiniGameControllerClickStatsFromDB(limit)
}

func loadRecentMonthlyMiniGameControllerClickStatsFromDB(limit int) []*entity.MonthlyMiniGameControllerClickStat {
	list := make([]*entity.MonthlyMiniGameControllerClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameControllerClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameControllerClickStats(list)
	return list
}

func reverseMonthlyMiniGameControllerClickStats(list []*entity.MonthlyMiniGameControllerClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameControllerClickStat(data *entity.MonthlyMiniGameControllerClickStat) {
	if data == nil || data.ID == "" || monthlyMiniGameControllerClickStatCacheMgr == nil {
		return
	}
	monthlyMiniGameControllerClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
