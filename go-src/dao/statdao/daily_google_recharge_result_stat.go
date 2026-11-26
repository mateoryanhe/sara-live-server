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
	dailyGoogleRechargeResultStatCacheMgr     *cache.RowCache[*entity.DailyGoogleRechargeResultStat]
)

func initDailyGoogleRechargeResultStatDao() {
	dailyGoogleRechargeResultStatCacheMgr = cache.NewRowCache[*entity.DailyGoogleRechargeResultStat]()
}

func GetDailyGoogleRechargeResultStatByDate(date string) *entity.DailyGoogleRechargeResultStat {
	return dailyGoogleRechargeResultStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGoogleRechargeResultStat, error) {
		var data *entity.DailyGoogleRechargeResultStat
		_ = g.Model(string(entity.TbDailyGoogleRechargeResultStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGoogleRechargeResultStat(date), nil
	})
}

func ListRecentDailyGoogleRechargeResultStats(limit int) []*entity.DailyGoogleRechargeResultStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyGoogleRechargeResultStatsFromDB(limit)
}

func loadRecentDailyGoogleRechargeResultStatsFromDB(limit int) []*entity.DailyGoogleRechargeResultStat {
	list := make([]*entity.DailyGoogleRechargeResultStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGoogleRechargeResultStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGoogleRechargeResultStats(list)
	return list
}

func reverseDailyGoogleRechargeResultStats(list []*entity.DailyGoogleRechargeResultStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGoogleRechargeResultStat(data *entity.DailyGoogleRechargeResultStat) {
	if data == nil || data.ID == "" || dailyGoogleRechargeResultStatCacheMgr == nil {
		return
	}
	dailyGoogleRechargeResultStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
