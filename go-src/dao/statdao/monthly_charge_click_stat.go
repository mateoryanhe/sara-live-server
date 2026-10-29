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
	monthlyChargeClickStatCacheMgr     *cache.RowCache[*entity.MonthlyChargeClickStat]
)

func initMonthlyChargeClickStatDao() {
	monthlyChargeClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyChargeClickStat]()
}

func GetMonthlyChargeClickStatByMonth(month string) *entity.MonthlyChargeClickStat {
	return monthlyChargeClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyChargeClickStat, error) {
		var data *entity.MonthlyChargeClickStat
		_ = g.Model(string(entity.TbMonthlyChargeClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyChargeClickStat(month), nil
	})
}

func ListRecentMonthlyChargeClickStats(limit int) []*entity.MonthlyChargeClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyChargeClickStatsFromDB(limit)
}

func loadRecentMonthlyChargeClickStatsFromDB(limit int) []*entity.MonthlyChargeClickStat {
	list := make([]*entity.MonthlyChargeClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyChargeClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyChargeClickStats(list)
	return list
}

func reverseMonthlyChargeClickStats(list []*entity.MonthlyChargeClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyChargeClickStat(data *entity.MonthlyChargeClickStat) {
	if data == nil || data.ID == "" || monthlyChargeClickStatCacheMgr == nil {
		return
	}
	monthlyChargeClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
