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
	weeklyMiniGameWebViewLoadStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameWebViewLoadStat]
	weeklyMiniGameWebViewLoadStatListCacheMgr *cache.ListCache[*entity.WeeklyMiniGameWebViewLoadStat]
)

func initWeeklyMiniGameWebViewLoadStatDao() {
	weeklyMiniGameWebViewLoadStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameWebViewLoadStat]()
	weeklyMiniGameWebViewLoadStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyMiniGameWebViewLoadStat]()
}

func GetWeeklyMiniGameWebViewLoadStatByWeek(week string) *entity.WeeklyMiniGameWebViewLoadStat {
	return weeklyMiniGameWebViewLoadStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameWebViewLoadStat, error) {
		var data *entity.WeeklyMiniGameWebViewLoadStat
		_ = g.Model(string(entity.TbWeeklyMiniGameWebViewLoadStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameWebViewLoadStat(week), nil
	})
}

func ListRecentWeeklyMiniGameWebViewLoadStats(limit int) []*entity.WeeklyMiniGameWebViewLoadStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyMiniGameWebViewLoadStatKey(time.Now()))
	return weeklyMiniGameWebViewLoadStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyMiniGameWebViewLoadStat, error) {
		return loadRecentWeeklyMiniGameWebViewLoadStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyMiniGameWebViewLoadStatsFromDB(limit int) []*entity.WeeklyMiniGameWebViewLoadStat {
	list := make([]*entity.WeeklyMiniGameWebViewLoadStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameWebViewLoadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameWebViewLoadStats(list)
	return list
}

func reverseWeeklyMiniGameWebViewLoadStats(list []*entity.WeeklyMiniGameWebViewLoadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameWebViewLoadStat(data *entity.WeeklyMiniGameWebViewLoadStat) {
	if data == nil || data.ID == "" || weeklyMiniGameWebViewLoadStatCacheMgr == nil {
		return
	}
	weeklyMiniGameWebViewLoadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
