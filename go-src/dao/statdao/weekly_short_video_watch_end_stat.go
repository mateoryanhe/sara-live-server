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
	weeklyShortVideoWatchEndStatCacheMgr     *cache.RowCache[*entity.WeeklyShortVideoWatchEndStat]
)

func initWeeklyShortVideoWatchEndStatDao() {
	weeklyShortVideoWatchEndStatCacheMgr = cache.NewRowCache[*entity.WeeklyShortVideoWatchEndStat]()
}

func GetWeeklyShortVideoWatchEndStatByWeek(week string) *entity.WeeklyShortVideoWatchEndStat {
	return weeklyShortVideoWatchEndStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyShortVideoWatchEndStat, error) {
		var data *entity.WeeklyShortVideoWatchEndStat
		_ = g.Model(string(entity.TbWeeklyShortVideoWatchEndStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyShortVideoWatchEndStat(week), nil
	})
}

func ListRecentWeeklyShortVideoWatchEndStats(limit int) []*entity.WeeklyShortVideoWatchEndStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyShortVideoWatchEndStatsFromDB(limit)
}

func loadRecentWeeklyShortVideoWatchEndStatsFromDB(limit int) []*entity.WeeklyShortVideoWatchEndStat {
	list := make([]*entity.WeeklyShortVideoWatchEndStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyShortVideoWatchEndStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyShortVideoWatchEndStats(list)
	return list
}

func reverseWeeklyShortVideoWatchEndStats(list []*entity.WeeklyShortVideoWatchEndStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyShortVideoWatchEndStat(data *entity.WeeklyShortVideoWatchEndStat) {
	if data == nil || data.ID == "" || weeklyShortVideoWatchEndStatCacheMgr == nil {
		return
	}
	weeklyShortVideoWatchEndStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
