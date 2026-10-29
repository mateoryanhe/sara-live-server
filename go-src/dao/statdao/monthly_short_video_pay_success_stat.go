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
	monthlyShortVideoPaySuccessStatCacheMgr     *cache.RowCache[*entity.MonthlyShortVideoPaySuccessStat]
)

func initMonthlyShortVideoPaySuccessStatDao() {
	monthlyShortVideoPaySuccessStatCacheMgr = cache.NewRowCache[*entity.MonthlyShortVideoPaySuccessStat]()
}

func GetMonthlyShortVideoPaySuccessStatByMonth(month string) *entity.MonthlyShortVideoPaySuccessStat {
	return monthlyShortVideoPaySuccessStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyShortVideoPaySuccessStat, error) {
		var data *entity.MonthlyShortVideoPaySuccessStat
		_ = g.Model(string(entity.TbMonthlyShortVideoPaySuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyShortVideoPaySuccessStat(month), nil
	})
}

func ListRecentMonthlyShortVideoPaySuccessStats(limit int) []*entity.MonthlyShortVideoPaySuccessStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyShortVideoPaySuccessStatsFromDB(limit)
}

func loadRecentMonthlyShortVideoPaySuccessStatsFromDB(limit int) []*entity.MonthlyShortVideoPaySuccessStat {
	list := make([]*entity.MonthlyShortVideoPaySuccessStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyShortVideoPaySuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyShortVideoPaySuccessStats(list)
	return list
}

func reverseMonthlyShortVideoPaySuccessStats(list []*entity.MonthlyShortVideoPaySuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyShortVideoPaySuccessStat(data *entity.MonthlyShortVideoPaySuccessStat) {
	if data == nil || data.ID == "" || monthlyShortVideoPaySuccessStatCacheMgr == nil {
		return
	}
	monthlyShortVideoPaySuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
