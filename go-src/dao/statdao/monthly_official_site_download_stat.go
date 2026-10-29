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
	monthlyOfficialSiteDownloadStatCacheMgr     *cache.RowCache[*entity.MonthlyOfficialSiteDownloadStat]
)

func initMonthlyOfficialSiteDownloadStatDao() {
	monthlyOfficialSiteDownloadStatCacheMgr = cache.NewRowCache[*entity.MonthlyOfficialSiteDownloadStat]()
}

func GetMonthlyOfficialSiteDownloadStatByMonth(month string) *entity.MonthlyOfficialSiteDownloadStat {
	return monthlyOfficialSiteDownloadStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyOfficialSiteDownloadStat, error) {
		var data *entity.MonthlyOfficialSiteDownloadStat
		_ = g.Model(string(entity.TbMonthlyOfficialSiteDownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyOfficialSiteDownloadStat(month), nil
	})
}

func ListRecentMonthlyOfficialSiteDownloadStats(limit int) []*entity.MonthlyOfficialSiteDownloadStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyOfficialSiteDownloadStatsFromDB(limit)
}

func loadRecentMonthlyOfficialSiteDownloadStatsFromDB(limit int) []*entity.MonthlyOfficialSiteDownloadStat {
	list := make([]*entity.MonthlyOfficialSiteDownloadStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyOfficialSiteDownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyOfficialSiteDownloadStats(list)
	return list
}

func reverseMonthlyOfficialSiteDownloadStats(list []*entity.MonthlyOfficialSiteDownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyOfficialSiteDownloadStat(data *entity.MonthlyOfficialSiteDownloadStat) {
	if data == nil || data.ID == "" || monthlyOfficialSiteDownloadStatCacheMgr == nil {
		return
	}
	monthlyOfficialSiteDownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
