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
	dailyGoogleRechargeListStatCacheMgr     *cache.RowCache[*entity.DailyGoogleRechargeListStat]
)

func initDailyGoogleRechargeListStatDao() {
	dailyGoogleRechargeListStatCacheMgr = cache.NewRowCache[*entity.DailyGoogleRechargeListStat]()
}

func GetDailyGoogleRechargeListStatByDate(date string) *entity.DailyGoogleRechargeListStat {
	return dailyGoogleRechargeListStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGoogleRechargeListStat, error) {
		var data *entity.DailyGoogleRechargeListStat
		_ = g.Model(string(entity.TbDailyGoogleRechargeListStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGoogleRechargeListStat(date), nil
	})
}

func ListRecentDailyGoogleRechargeListStats(limit int) []*entity.DailyGoogleRechargeListStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyGoogleRechargeListStatsFromDB(limit)
}

func loadRecentDailyGoogleRechargeListStatsFromDB(limit int) []*entity.DailyGoogleRechargeListStat {
	list := make([]*entity.DailyGoogleRechargeListStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGoogleRechargeListStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGoogleRechargeListStats(list)
	return list
}

func reverseDailyGoogleRechargeListStats(list []*entity.DailyGoogleRechargeListStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGoogleRechargeListStat(data *entity.DailyGoogleRechargeListStat) {
	if data == nil || data.ID == "" || dailyGoogleRechargeListStatCacheMgr == nil {
		return
	}
	dailyGoogleRechargeListStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
