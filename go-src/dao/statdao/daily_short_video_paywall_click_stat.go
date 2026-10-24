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
	dailyShortVideoPaywallClickStatCacheMgr     *cache.RowCache[*entity.DailyShortVideoPaywallClickStat]
)

func initDailyShortVideoPaywallClickStatDao() {
	dailyShortVideoPaywallClickStatCacheMgr = cache.NewRowCache[*entity.DailyShortVideoPaywallClickStat]()
}

func GetDailyShortVideoPaywallClickStatByDate(date string) *entity.DailyShortVideoPaywallClickStat {
	return dailyShortVideoPaywallClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyShortVideoPaywallClickStat, error) {
		var data *entity.DailyShortVideoPaywallClickStat
		_ = g.Model(string(entity.TbDailyShortVideoPaywallClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyShortVideoPaywallClickStat(date), nil
	})
}

func ListRecentDailyShortVideoPaywallClickStats(limit int) []*entity.DailyShortVideoPaywallClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyShortVideoPaywallClickStatsFromDB(limit)
}

func loadRecentDailyShortVideoPaywallClickStatsFromDB(limit int) []*entity.DailyShortVideoPaywallClickStat {
	list := make([]*entity.DailyShortVideoPaywallClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyShortVideoPaywallClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyShortVideoPaywallClickStats(list)
	return list
}

func reverseDailyShortVideoPaywallClickStats(list []*entity.DailyShortVideoPaywallClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyShortVideoPaywallClickStat(data *entity.DailyShortVideoPaywallClickStat) {
	if data == nil || data.ID == "" || dailyShortVideoPaywallClickStatCacheMgr == nil {
		return
	}
	dailyShortVideoPaywallClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
