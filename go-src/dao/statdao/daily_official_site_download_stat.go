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
	dailyOfficialSiteDownloadStatCacheMgr *cache.RowCache[*entity.DailyOfficialSiteDownloadStat]
)

func initDailyOfficialSiteDownloadStatDao() {
	dailyOfficialSiteDownloadStatCacheMgr = cache.NewRowCache[*entity.DailyOfficialSiteDownloadStat]()
}

func GetDailyOfficialSiteDownloadStatByDate(date string) *entity.DailyOfficialSiteDownloadStat {
	return dailyOfficialSiteDownloadStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyOfficialSiteDownloadStat, error) {
		var data *entity.DailyOfficialSiteDownloadStat
		_ = g.Model(string(entity.TbDailyOfficialSiteDownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyOfficialSiteDownloadStat(date), nil
	})
}

func ListRecentDailyOfficialSiteDownloadStats(limit int) []*entity.DailyOfficialSiteDownloadStat {
	if limit <= 0 {
		limit = 30
	}
	return loadRecentDailyOfficialSiteDownloadStatsFromDB(limit)
}

func loadRecentDailyOfficialSiteDownloadStatsFromDB(limit int) []*entity.DailyOfficialSiteDownloadStat {
	list := make([]*entity.DailyOfficialSiteDownloadStat, 0, limit)
	_ = g.Model(string(entity.TbDailyOfficialSiteDownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyOfficialSiteDownloadStats(list)
	return list
}

func reverseDailyOfficialSiteDownloadStats(list []*entity.DailyOfficialSiteDownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyOfficialSiteDownloadStat(data *entity.DailyOfficialSiteDownloadStat) {
	if data == nil || data.ID == "" || dailyOfficialSiteDownloadStatCacheMgr == nil {
		return
	}
	dailyOfficialSiteDownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
