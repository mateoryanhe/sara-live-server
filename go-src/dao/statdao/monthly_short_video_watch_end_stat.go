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
	monthlyShortVideoWatchEndStatCacheMgr     *cache.RowCache[*entity.MonthlyShortVideoWatchEndStat]
)

func initMonthlyShortVideoWatchEndStatDao() {
	monthlyShortVideoWatchEndStatCacheMgr = cache.NewRowCache[*entity.MonthlyShortVideoWatchEndStat]()
}

func GetMonthlyShortVideoWatchEndStatByMonth(month string) *entity.MonthlyShortVideoWatchEndStat {
	return monthlyShortVideoWatchEndStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyShortVideoWatchEndStat, error) {
		var data *entity.MonthlyShortVideoWatchEndStat
		_ = g.Model(string(entity.TbMonthlyShortVideoWatchEndStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyShortVideoWatchEndStat(month), nil
	})
}

func ListRecentMonthlyShortVideoWatchEndStats(limit int) []*entity.MonthlyShortVideoWatchEndStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyShortVideoWatchEndStatsFromDB(limit)
}

func loadRecentMonthlyShortVideoWatchEndStatsFromDB(limit int) []*entity.MonthlyShortVideoWatchEndStat {
	list := make([]*entity.MonthlyShortVideoWatchEndStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyShortVideoWatchEndStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyShortVideoWatchEndStats(list)
	return list
}

func reverseMonthlyShortVideoWatchEndStats(list []*entity.MonthlyShortVideoWatchEndStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyShortVideoWatchEndStat(data *entity.MonthlyShortVideoWatchEndStat) {
	if data == nil || data.ID == "" || monthlyShortVideoWatchEndStatCacheMgr == nil {
		return
	}
	monthlyShortVideoWatchEndStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
