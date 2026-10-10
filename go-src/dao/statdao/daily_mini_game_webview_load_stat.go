package statdao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var (
	dailyMiniGameWebViewLoadStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameWebViewLoadStat]
	dailyMiniGameWebViewLoadStatListCacheMgr *cache.ListCache[*entity.DailyMiniGameWebViewLoadStat]
)

func initDailyMiniGameWebViewLoadStatDao() {
	dailyMiniGameWebViewLoadStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameWebViewLoadStat]()
	dailyMiniGameWebViewLoadStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyMiniGameWebViewLoadStat]()
}

func GetDailyMiniGameWebViewLoadStatByDate(date string) *entity.DailyMiniGameWebViewLoadStat {
	return dailyMiniGameWebViewLoadStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameWebViewLoadStat, error) {
		var data *entity.DailyMiniGameWebViewLoadStat
		_ = g.Model(string(entity.TbDailyMiniGameWebViewLoadStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameWebViewLoadStat(date), nil
	})
}

func ListRecentDailyMiniGameWebViewLoadStats(limit int) []*entity.DailyMiniGameWebViewLoadStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyMiniGameWebViewLoadStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyMiniGameWebViewLoadStat, error) {
		return loadRecentDailyMiniGameWebViewLoadStatsFromDB(limit), nil
	})
}

func loadRecentDailyMiniGameWebViewLoadStatsFromDB(limit int) []*entity.DailyMiniGameWebViewLoadStat {
	list := make([]*entity.DailyMiniGameWebViewLoadStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameWebViewLoadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameWebViewLoadStats(list)
	return list
}

func reverseDailyMiniGameWebViewLoadStats(list []*entity.DailyMiniGameWebViewLoadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameWebViewLoadStat(data *entity.DailyMiniGameWebViewLoadStat) {
	if data == nil || data.ID == "" || dailyMiniGameWebViewLoadStatCacheMgr == nil {
		return
	}
	dailyMiniGameWebViewLoadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
