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
	dailyH5DownloadStatCacheMgr *cache.RowCache[*entity.DailyH5DownloadStat]
)

func initDailyH5DownloadStatDao() {
	dailyH5DownloadStatCacheMgr = cache.NewRowCache[*entity.DailyH5DownloadStat]()
}

func GetDailyH5DownloadStatByDate(date string) *entity.DailyH5DownloadStat {
	return dailyH5DownloadStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyH5DownloadStat, error) {
		var data *entity.DailyH5DownloadStat
		_ = g.Model(string(entity.TbDailyH5DownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyH5DownloadStat(date), nil
	})
}

func ListRecentDailyH5DownloadStats(limit int) []*entity.DailyH5DownloadStat {
	if limit <= 0 {
		limit = 30
	}
	return loadRecentDailyH5DownloadStatsFromDB(limit)
}

func loadRecentDailyH5DownloadStatsFromDB(limit int) []*entity.DailyH5DownloadStat {
	list := make([]*entity.DailyH5DownloadStat, 0, limit)
	_ = g.Model(string(entity.TbDailyH5DownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyH5DownloadStats(list)
	return list
}

func reverseDailyH5DownloadStats(list []*entity.DailyH5DownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyH5DownloadStat(data *entity.DailyH5DownloadStat) {
	if data == nil || data.ID == "" || dailyH5DownloadStatCacheMgr == nil {
		return
	}
	dailyH5DownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
