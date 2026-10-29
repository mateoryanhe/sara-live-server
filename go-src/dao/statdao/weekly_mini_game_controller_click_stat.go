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
	weeklyMiniGameControllerClickStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameControllerClickStat]
)

func initWeeklyMiniGameControllerClickStatDao() {
	weeklyMiniGameControllerClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameControllerClickStat]()
}

func GetWeeklyMiniGameControllerClickStatByWeek(week string) *entity.WeeklyMiniGameControllerClickStat {
	return weeklyMiniGameControllerClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameControllerClickStat, error) {
		var data *entity.WeeklyMiniGameControllerClickStat
		_ = g.Model(string(entity.TbWeeklyMiniGameControllerClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameControllerClickStat(week), nil
	})
}

func ListRecentWeeklyMiniGameControllerClickStats(limit int) []*entity.WeeklyMiniGameControllerClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyMiniGameControllerClickStatsFromDB(limit)
}

func loadRecentWeeklyMiniGameControllerClickStatsFromDB(limit int) []*entity.WeeklyMiniGameControllerClickStat {
	list := make([]*entity.WeeklyMiniGameControllerClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameControllerClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameControllerClickStats(list)
	return list
}

func reverseWeeklyMiniGameControllerClickStats(list []*entity.WeeklyMiniGameControllerClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameControllerClickStat(data *entity.WeeklyMiniGameControllerClickStat) {
	if data == nil || data.ID == "" || weeklyMiniGameControllerClickStatCacheMgr == nil {
		return
	}
	weeklyMiniGameControllerClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
