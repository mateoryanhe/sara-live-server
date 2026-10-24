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
	dailyGameLiveWatchStreamStatCacheMgr     *cache.RowCache[*entity.DailyGameLiveWatchStreamStat]
)

func initDailyGameLiveWatchStreamStatDao() {
	dailyGameLiveWatchStreamStatCacheMgr = cache.NewRowCache[*entity.DailyGameLiveWatchStreamStat]()
}

func GetDailyGameLiveWatchStreamStatByDate(date string) *entity.DailyGameLiveWatchStreamStat {
	return dailyGameLiveWatchStreamStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGameLiveWatchStreamStat, error) {
		var data *entity.DailyGameLiveWatchStreamStat
		_ = g.Model(string(entity.TbDailyGameLiveWatchStreamStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGameLiveWatchStreamStat(date), nil
	})
}

func ListRecentDailyGameLiveWatchStreamStats(limit int) []*entity.DailyGameLiveWatchStreamStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyGameLiveWatchStreamStatsFromDB(limit)
}

func loadRecentDailyGameLiveWatchStreamStatsFromDB(limit int) []*entity.DailyGameLiveWatchStreamStat {
	list := make([]*entity.DailyGameLiveWatchStreamStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGameLiveWatchStreamStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGameLiveWatchStreamStats(list)
	return list
}

func reverseDailyGameLiveWatchStreamStats(list []*entity.DailyGameLiveWatchStreamStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGameLiveWatchStreamStat(data *entity.DailyGameLiveWatchStreamStat) {
	if data == nil || data.ID == "" || dailyGameLiveWatchStreamStatCacheMgr == nil {
		return
	}
	dailyGameLiveWatchStreamStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
