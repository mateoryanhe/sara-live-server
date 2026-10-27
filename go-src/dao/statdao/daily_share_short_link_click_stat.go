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
	dailyShareShortLinkClickStatCacheMgr *cache.RowCache[*entity.DailyShareShortLinkClickStat]
)

func initDailyShareShortLinkClickStatDao() {
	dailyShareShortLinkClickStatCacheMgr = cache.NewRowCache[*entity.DailyShareShortLinkClickStat]()
}

func GetDailyShareShortLinkClickStatByDate(date string) *entity.DailyShareShortLinkClickStat {
	return dailyShareShortLinkClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyShareShortLinkClickStat, error) {
		var data *entity.DailyShareShortLinkClickStat
		_ = g.Model(string(entity.TbDailyShareShortLinkClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyShareShortLinkClickStat(date), nil
	})
}

func ListRecentDailyShareShortLinkClickStats(limit int) []*entity.DailyShareShortLinkClickStat {
	if limit <= 0 {
		limit = 30
	}
	return loadRecentDailyShareShortLinkClickStatsFromDB(limit)
}

func loadRecentDailyShareShortLinkClickStatsFromDB(limit int) []*entity.DailyShareShortLinkClickStat {
	list := make([]*entity.DailyShareShortLinkClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyShareShortLinkClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyShareShortLinkClickStats(list)
	return list
}

func reverseDailyShareShortLinkClickStats(list []*entity.DailyShareShortLinkClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyShareShortLinkClickStat(data *entity.DailyShareShortLinkClickStat) {
	if data == nil || data.ID == "" || dailyShareShortLinkClickStatCacheMgr == nil {
		return
	}
	dailyShareShortLinkClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
