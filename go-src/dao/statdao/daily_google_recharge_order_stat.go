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
	dailyGoogleRechargeOrderStatCacheMgr     *cache.RowCache[*entity.DailyGoogleRechargeOrderStat]
)

func initDailyGoogleRechargeOrderStatDao() {
	dailyGoogleRechargeOrderStatCacheMgr = cache.NewRowCache[*entity.DailyGoogleRechargeOrderStat]()
}

func GetDailyGoogleRechargeOrderStatByDate(date string) *entity.DailyGoogleRechargeOrderStat {
	return dailyGoogleRechargeOrderStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGoogleRechargeOrderStat, error) {
		var data *entity.DailyGoogleRechargeOrderStat
		_ = g.Model(string(entity.TbDailyGoogleRechargeOrderStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGoogleRechargeOrderStat(date), nil
	})
}

func ListRecentDailyGoogleRechargeOrderStats(limit int) []*entity.DailyGoogleRechargeOrderStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyGoogleRechargeOrderStatsFromDB(limit)
}

func loadRecentDailyGoogleRechargeOrderStatsFromDB(limit int) []*entity.DailyGoogleRechargeOrderStat {
	list := make([]*entity.DailyGoogleRechargeOrderStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGoogleRechargeOrderStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGoogleRechargeOrderStats(list)
	return list
}

func reverseDailyGoogleRechargeOrderStats(list []*entity.DailyGoogleRechargeOrderStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGoogleRechargeOrderStat(data *entity.DailyGoogleRechargeOrderStat) {
	if data == nil || data.ID == "" || dailyGoogleRechargeOrderStatCacheMgr == nil {
		return
	}
	dailyGoogleRechargeOrderStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
