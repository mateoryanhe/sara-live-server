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
	monthlyH5DownloadStatCacheMgr     *cache.RowCache[*entity.MonthlyH5DownloadStat]
)

func initMonthlyH5DownloadStatDao() {
	monthlyH5DownloadStatCacheMgr = cache.NewRowCache[*entity.MonthlyH5DownloadStat]()
}

func GetMonthlyH5DownloadStatByMonth(month string) *entity.MonthlyH5DownloadStat {
	return monthlyH5DownloadStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyH5DownloadStat, error) {
		var data *entity.MonthlyH5DownloadStat
		_ = g.Model(string(entity.TbMonthlyH5DownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyH5DownloadStat(month), nil
	})
}

func ListRecentMonthlyH5DownloadStats(limit int) []*entity.MonthlyH5DownloadStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyH5DownloadStatsFromDB(limit)
}

func loadRecentMonthlyH5DownloadStatsFromDB(limit int) []*entity.MonthlyH5DownloadStat {
	list := make([]*entity.MonthlyH5DownloadStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyH5DownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyH5DownloadStats(list)
	return list
}

func reverseMonthlyH5DownloadStats(list []*entity.MonthlyH5DownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyH5DownloadStat(data *entity.MonthlyH5DownloadStat) {
	if data == nil || data.ID == "" || monthlyH5DownloadStatCacheMgr == nil {
		return
	}
	monthlyH5DownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
