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
	monthlyShortVideoPaywallClickStatCacheMgr     *cache.RowCache[*entity.MonthlyShortVideoPaywallClickStat]
)

func initMonthlyShortVideoPaywallClickStatDao() {
	monthlyShortVideoPaywallClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyShortVideoPaywallClickStat]()
}

func GetMonthlyShortVideoPaywallClickStatByMonth(month string) *entity.MonthlyShortVideoPaywallClickStat {
	return monthlyShortVideoPaywallClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyShortVideoPaywallClickStat, error) {
		var data *entity.MonthlyShortVideoPaywallClickStat
		_ = g.Model(string(entity.TbMonthlyShortVideoPaywallClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyShortVideoPaywallClickStat(month), nil
	})
}

func ListRecentMonthlyShortVideoPaywallClickStats(limit int) []*entity.MonthlyShortVideoPaywallClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyShortVideoPaywallClickStatsFromDB(limit)
}

func loadRecentMonthlyShortVideoPaywallClickStatsFromDB(limit int) []*entity.MonthlyShortVideoPaywallClickStat {
	list := make([]*entity.MonthlyShortVideoPaywallClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyShortVideoPaywallClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyShortVideoPaywallClickStats(list)
	return list
}

func reverseMonthlyShortVideoPaywallClickStats(list []*entity.MonthlyShortVideoPaywallClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyShortVideoPaywallClickStat(data *entity.MonthlyShortVideoPaywallClickStat) {
	if data == nil || data.ID == "" || monthlyShortVideoPaywallClickStatCacheMgr == nil {
		return
	}
	monthlyShortVideoPaywallClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
