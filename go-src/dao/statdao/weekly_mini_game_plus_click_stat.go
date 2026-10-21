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
	weeklyMiniGamePlusClickStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGamePlusClickStat]
	weeklyMiniGamePlusClickStatListCacheMgr *cache.ListCache[*entity.WeeklyMiniGamePlusClickStat]
)

func initWeeklyMiniGamePlusClickStatDao() {
	weeklyMiniGamePlusClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGamePlusClickStat]()
	weeklyMiniGamePlusClickStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyMiniGamePlusClickStat]()
}

func GetWeeklyMiniGamePlusClickStatByWeek(week string) *entity.WeeklyMiniGamePlusClickStat {
	return weeklyMiniGamePlusClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGamePlusClickStat, error) {
		var data *entity.WeeklyMiniGamePlusClickStat
		_ = g.Model(string(entity.TbWeeklyMiniGamePlusClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGamePlusClickStat(week), nil
	})
}

func ListRecentWeeklyMiniGamePlusClickStats(limit int) []*entity.WeeklyMiniGamePlusClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyMiniGamePlusClickStatKey(time.Now()))
	return weeklyMiniGamePlusClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyMiniGamePlusClickStat, error) {
		return loadRecentWeeklyMiniGamePlusClickStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyMiniGamePlusClickStatsFromDB(limit int) []*entity.WeeklyMiniGamePlusClickStat {
	list := make([]*entity.WeeklyMiniGamePlusClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGamePlusClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGamePlusClickStats(list)
	return list
}

func reverseWeeklyMiniGamePlusClickStats(list []*entity.WeeklyMiniGamePlusClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGamePlusClickStat(data *entity.WeeklyMiniGamePlusClickStat) {
	if data == nil || data.ID == "" || weeklyMiniGamePlusClickStatCacheMgr == nil {
		return
	}
	weeklyMiniGamePlusClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
