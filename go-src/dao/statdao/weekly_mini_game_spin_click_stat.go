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
	weeklyMiniGameSpinClickStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameSpinClickStat]
)

func initWeeklyMiniGameSpinClickStatDao() {
	weeklyMiniGameSpinClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameSpinClickStat]()
}

func GetWeeklyMiniGameSpinClickStatByWeek(week string) *entity.WeeklyMiniGameSpinClickStat {
	return weeklyMiniGameSpinClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameSpinClickStat, error) {
		var data *entity.WeeklyMiniGameSpinClickStat
		_ = g.Model(string(entity.TbWeeklyMiniGameSpinClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameSpinClickStat(week), nil
	})
}

func ListRecentWeeklyMiniGameSpinClickStats(limit int) []*entity.WeeklyMiniGameSpinClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyMiniGameSpinClickStatsFromDB(limit)
}

func loadRecentWeeklyMiniGameSpinClickStatsFromDB(limit int) []*entity.WeeklyMiniGameSpinClickStat {
	list := make([]*entity.WeeklyMiniGameSpinClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameSpinClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameSpinClickStats(list)
	return list
}

func reverseWeeklyMiniGameSpinClickStats(list []*entity.WeeklyMiniGameSpinClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameSpinClickStat(data *entity.WeeklyMiniGameSpinClickStat) {
	if data == nil || data.ID == "" || weeklyMiniGameSpinClickStatCacheMgr == nil {
		return
	}
	weeklyMiniGameSpinClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
