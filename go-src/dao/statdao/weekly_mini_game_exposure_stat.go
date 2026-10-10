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
	weeklyMiniGameExposureStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameExposureStat]
	weeklyMiniGameExposureStatListCacheMgr *cache.ListCache[*entity.WeeklyMiniGameExposureStat]
)

func initWeeklyMiniGameExposureStatDao() {
	weeklyMiniGameExposureStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameExposureStat]()
	weeklyMiniGameExposureStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyMiniGameExposureStat]()
}

func GetWeeklyMiniGameExposureStatByWeek(week string) *entity.WeeklyMiniGameExposureStat {
	return weeklyMiniGameExposureStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameExposureStat, error) {
		var data *entity.WeeklyMiniGameExposureStat
		_ = g.Model(string(entity.TbWeeklyMiniGameExposureStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameExposureStat(week), nil
	})
}

func ListRecentWeeklyMiniGameExposureStats(limit int) []*entity.WeeklyMiniGameExposureStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyMiniGameExposureStatKey(time.Now()))
	return weeklyMiniGameExposureStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyMiniGameExposureStat, error) {
		return loadRecentWeeklyMiniGameExposureStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyMiniGameExposureStatsFromDB(limit int) []*entity.WeeklyMiniGameExposureStat {
	list := make([]*entity.WeeklyMiniGameExposureStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameExposureStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameExposureStats(list)
	return list
}

func reverseWeeklyMiniGameExposureStats(list []*entity.WeeklyMiniGameExposureStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameExposureStat(data *entity.WeeklyMiniGameExposureStat) {
	if data == nil || data.ID == "" || weeklyMiniGameExposureStatCacheMgr == nil {
		return
	}
	weeklyMiniGameExposureStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
