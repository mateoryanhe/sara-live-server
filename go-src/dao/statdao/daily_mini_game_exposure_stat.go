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
	dailyMiniGameExposureStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameExposureStat]
	dailyMiniGameExposureStatListCacheMgr *cache.ListCache[*entity.DailyMiniGameExposureStat]
)

func initDailyMiniGameExposureStatDao() {
	dailyMiniGameExposureStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameExposureStat]()
	dailyMiniGameExposureStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyMiniGameExposureStat]()
}

func GetDailyMiniGameExposureStatByDate(date string) *entity.DailyMiniGameExposureStat {
	return dailyMiniGameExposureStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameExposureStat, error) {
		var data *entity.DailyMiniGameExposureStat
		_ = g.Model(string(entity.TbDailyMiniGameExposureStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameExposureStat(date), nil
	})
}

func ListRecentDailyMiniGameExposureStats(limit int) []*entity.DailyMiniGameExposureStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyMiniGameExposureStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyMiniGameExposureStat, error) {
		return loadRecentDailyMiniGameExposureStatsFromDB(limit), nil
	})
}

func loadRecentDailyMiniGameExposureStatsFromDB(limit int) []*entity.DailyMiniGameExposureStat {
	list := make([]*entity.DailyMiniGameExposureStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameExposureStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameExposureStats(list)
	return list
}

func reverseDailyMiniGameExposureStats(list []*entity.DailyMiniGameExposureStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameExposureStat(data *entity.DailyMiniGameExposureStat) {
	if data == nil || data.ID == "" || dailyMiniGameExposureStatCacheMgr == nil {
		return
	}
	dailyMiniGameExposureStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
