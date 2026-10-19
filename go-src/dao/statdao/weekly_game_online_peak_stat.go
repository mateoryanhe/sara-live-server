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
	weeklyGameOnlinePeakStatCacheMgr     *cache.RowCache[*entity.WeeklyGameOnlinePeakStat]
	weeklyGameOnlinePeakStatListCacheMgr *cache.ListCache[*entity.WeeklyGameOnlinePeakStat]
)

func initWeeklyGameOnlinePeakStatDao() {
	weeklyGameOnlinePeakStatCacheMgr = cache.NewRowCache[*entity.WeeklyGameOnlinePeakStat]()
	weeklyGameOnlinePeakStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyGameOnlinePeakStat]()
}

func GetWeeklyGameOnlinePeakStatByWeek(week string) *entity.WeeklyGameOnlinePeakStat {
	return weeklyGameOnlinePeakStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyGameOnlinePeakStat, error) {
		var data *entity.WeeklyGameOnlinePeakStat
		_ = g.Model(string(entity.TbWeeklyGameOnlinePeakStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyGameOnlinePeakStat(week), nil
	})
}

func ListRecentWeeklyGameOnlinePeakStats(limit int) []*entity.WeeklyGameOnlinePeakStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyLoginStatKey(time.Now()))
	return weeklyGameOnlinePeakStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyGameOnlinePeakStat, error) {
		return loadRecentWeeklyGameOnlinePeakStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyGameOnlinePeakStatsFromDB(limit int) []*entity.WeeklyGameOnlinePeakStat {
	list := make([]*entity.WeeklyGameOnlinePeakStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyGameOnlinePeakStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyGameOnlinePeakStats(list)
	return list
}

func reverseWeeklyGameOnlinePeakStats(list []*entity.WeeklyGameOnlinePeakStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyGameOnlinePeakStat(data *entity.WeeklyGameOnlinePeakStat) {
	if data == nil || data.ID == "" || weeklyGameOnlinePeakStatCacheMgr == nil {
		return
	}
	weeklyGameOnlinePeakStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
