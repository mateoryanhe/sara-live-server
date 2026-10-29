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
	weeklyChargePanelEnterStatCacheMgr     *cache.RowCache[*entity.WeeklyChargePanelEnterStat]
)

func initWeeklyChargePanelEnterStatDao() {
	weeklyChargePanelEnterStatCacheMgr = cache.NewRowCache[*entity.WeeklyChargePanelEnterStat]()
}

func GetWeeklyChargePanelEnterStatByWeek(week string) *entity.WeeklyChargePanelEnterStat {
	return weeklyChargePanelEnterStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyChargePanelEnterStat, error) {
		var data *entity.WeeklyChargePanelEnterStat
		_ = g.Model(string(entity.TbWeeklyChargePanelEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyChargePanelEnterStat(week), nil
	})
}

func ListRecentWeeklyChargePanelEnterStats(limit int) []*entity.WeeklyChargePanelEnterStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyChargePanelEnterStatsFromDB(limit)
}

func loadRecentWeeklyChargePanelEnterStatsFromDB(limit int) []*entity.WeeklyChargePanelEnterStat {
	list := make([]*entity.WeeklyChargePanelEnterStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyChargePanelEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyChargePanelEnterStats(list)
	return list
}

func reverseWeeklyChargePanelEnterStats(list []*entity.WeeklyChargePanelEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyChargePanelEnterStat(data *entity.WeeklyChargePanelEnterStat) {
	if data == nil || data.ID == "" || weeklyChargePanelEnterStatCacheMgr == nil {
		return
	}
	weeklyChargePanelEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
