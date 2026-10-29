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
	monthlyMiniGameSpinClickStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameSpinClickStat]
)

func initMonthlyMiniGameSpinClickStatDao() {
	monthlyMiniGameSpinClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameSpinClickStat]()
}

func GetMonthlyMiniGameSpinClickStatByMonth(month string) *entity.MonthlyMiniGameSpinClickStat {
	return monthlyMiniGameSpinClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameSpinClickStat, error) {
		var data *entity.MonthlyMiniGameSpinClickStat
		_ = g.Model(string(entity.TbMonthlyMiniGameSpinClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameSpinClickStat(month), nil
	})
}

func ListRecentMonthlyMiniGameSpinClickStats(limit int) []*entity.MonthlyMiniGameSpinClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyMiniGameSpinClickStatsFromDB(limit)
}

func loadRecentMonthlyMiniGameSpinClickStatsFromDB(limit int) []*entity.MonthlyMiniGameSpinClickStat {
	list := make([]*entity.MonthlyMiniGameSpinClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameSpinClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameSpinClickStats(list)
	return list
}

func reverseMonthlyMiniGameSpinClickStats(list []*entity.MonthlyMiniGameSpinClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameSpinClickStat(data *entity.MonthlyMiniGameSpinClickStat) {
	if data == nil || data.ID == "" || monthlyMiniGameSpinClickStatCacheMgr == nil {
		return
	}
	monthlyMiniGameSpinClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
