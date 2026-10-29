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
	weeklyShortVideoPlayStartStatCacheMgr     *cache.RowCache[*entity.WeeklyShortVideoPlayStartStat]
)

func initWeeklyShortVideoPlayStartStatDao() {
	weeklyShortVideoPlayStartStatCacheMgr = cache.NewRowCache[*entity.WeeklyShortVideoPlayStartStat]()
}

func GetWeeklyShortVideoPlayStartStatByWeek(week string) *entity.WeeklyShortVideoPlayStartStat {
	return weeklyShortVideoPlayStartStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyShortVideoPlayStartStat, error) {
		var data *entity.WeeklyShortVideoPlayStartStat
		_ = g.Model(string(entity.TbWeeklyShortVideoPlayStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyShortVideoPlayStartStat(week), nil
	})
}

func ListRecentWeeklyShortVideoPlayStartStats(limit int) []*entity.WeeklyShortVideoPlayStartStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyShortVideoPlayStartStatsFromDB(limit)
}

func loadRecentWeeklyShortVideoPlayStartStatsFromDB(limit int) []*entity.WeeklyShortVideoPlayStartStat {
	list := make([]*entity.WeeklyShortVideoPlayStartStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyShortVideoPlayStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyShortVideoPlayStartStats(list)
	return list
}

func reverseWeeklyShortVideoPlayStartStats(list []*entity.WeeklyShortVideoPlayStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyShortVideoPlayStartStat(data *entity.WeeklyShortVideoPlayStartStat) {
	if data == nil || data.ID == "" || weeklyShortVideoPlayStartStatCacheMgr == nil {
		return
	}
	weeklyShortVideoPlayStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
