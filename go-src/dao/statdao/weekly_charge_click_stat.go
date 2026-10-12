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
	weeklyChargeClickStatCacheMgr     *cache.RowCache[*entity.WeeklyChargeClickStat]
	weeklyChargeClickStatListCacheMgr *cache.ListCache[*entity.WeeklyChargeClickStat]
)

func initWeeklyChargeClickStatDao() {
	weeklyChargeClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyChargeClickStat]()
	weeklyChargeClickStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyChargeClickStat]()
}

func GetWeeklyChargeClickStatByWeek(week string) *entity.WeeklyChargeClickStat {
	return weeklyChargeClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyChargeClickStat, error) {
		var data *entity.WeeklyChargeClickStat
		_ = g.Model(string(entity.TbWeeklyChargeClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyChargeClickStat(week), nil
	})
}

func ListRecentWeeklyChargeClickStats(limit int) []*entity.WeeklyChargeClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyChargeClickStatKey(time.Now()))
	return weeklyChargeClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyChargeClickStat, error) {
		return loadRecentWeeklyChargeClickStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyChargeClickStatsFromDB(limit int) []*entity.WeeklyChargeClickStat {
	list := make([]*entity.WeeklyChargeClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyChargeClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyChargeClickStats(list)
	return list
}

func reverseWeeklyChargeClickStats(list []*entity.WeeklyChargeClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyChargeClickStat(data *entity.WeeklyChargeClickStat) {
	if data == nil || data.ID == "" || weeklyChargeClickStatCacheMgr == nil {
		return
	}
	weeklyChargeClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
