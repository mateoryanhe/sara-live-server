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
	monthlyChargeResultServerStatCacheMgr     *cache.RowCache[*entity.MonthlyChargeResultServerStat]
)

func initMonthlyChargeResultServerStatDao() {
	monthlyChargeResultServerStatCacheMgr = cache.NewRowCache[*entity.MonthlyChargeResultServerStat]()
}

func GetMonthlyChargeResultServerStatByMonth(month string) *entity.MonthlyChargeResultServerStat {
	return monthlyChargeResultServerStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyChargeResultServerStat, error) {
		var data *entity.MonthlyChargeResultServerStat
		_ = g.Model(string(entity.TbMonthlyChargeResultServerStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyChargeResultServerStat(month), nil
	})
}

func ListRecentMonthlyChargeResultServerStats(limit int) []*entity.MonthlyChargeResultServerStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyChargeResultServerStatsFromDB(limit)
}

func loadRecentMonthlyChargeResultServerStatsFromDB(limit int) []*entity.MonthlyChargeResultServerStat {
	list := make([]*entity.MonthlyChargeResultServerStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyChargeResultServerStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyChargeResultServerStats(list)
	return list
}

func reverseMonthlyChargeResultServerStats(list []*entity.MonthlyChargeResultServerStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyChargeResultServerStat(data *entity.MonthlyChargeResultServerStat) {
	if data == nil || data.ID == "" || monthlyChargeResultServerStatCacheMgr == nil {
		return
	}
	monthlyChargeResultServerStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
