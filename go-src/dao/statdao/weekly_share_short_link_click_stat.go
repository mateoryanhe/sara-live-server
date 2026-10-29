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
	weeklyShareShortLinkClickStatCacheMgr     *cache.RowCache[*entity.WeeklyShareShortLinkClickStat]
)

func initWeeklyShareShortLinkClickStatDao() {
	weeklyShareShortLinkClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyShareShortLinkClickStat]()
}

func GetWeeklyShareShortLinkClickStatByWeek(week string) *entity.WeeklyShareShortLinkClickStat {
	return weeklyShareShortLinkClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyShareShortLinkClickStat, error) {
		var data *entity.WeeklyShareShortLinkClickStat
		_ = g.Model(string(entity.TbWeeklyShareShortLinkClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyShareShortLinkClickStat(week), nil
	})
}

func ListRecentWeeklyShareShortLinkClickStats(limit int) []*entity.WeeklyShareShortLinkClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyShareShortLinkClickStatsFromDB(limit)
}

func loadRecentWeeklyShareShortLinkClickStatsFromDB(limit int) []*entity.WeeklyShareShortLinkClickStat {
	list := make([]*entity.WeeklyShareShortLinkClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyShareShortLinkClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyShareShortLinkClickStats(list)
	return list
}

func reverseWeeklyShareShortLinkClickStats(list []*entity.WeeklyShareShortLinkClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyShareShortLinkClickStat(data *entity.WeeklyShareShortLinkClickStat) {
	if data == nil || data.ID == "" || weeklyShareShortLinkClickStatCacheMgr == nil {
		return
	}
	weeklyShareShortLinkClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
