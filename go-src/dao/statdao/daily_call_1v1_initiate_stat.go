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
	dailyCall1v1InitiateStatCacheMgr *cache.RowCache[*entity.DailyCall1v1InitiateStat]
)

func initDailyCall1v1InitiateStatDao() {
	dailyCall1v1InitiateStatCacheMgr = cache.NewRowCache[*entity.DailyCall1v1InitiateStat]()
}

func GetDailyCall1v1InitiateStatByDate(date string) *entity.DailyCall1v1InitiateStat {
	return dailyCall1v1InitiateStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyCall1v1InitiateStat, error) {
		var data *entity.DailyCall1v1InitiateStat
		_ = g.Model(string(entity.TbDailyCall1v1InitiateStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyCall1v1InitiateStat(date), nil
	})
}

func ListRecentDailyCall1v1InitiateStats(limit int) []*entity.DailyCall1v1InitiateStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyCall1v1InitiateStatsFromDB(limit)
}

func loadRecentDailyCall1v1InitiateStatsFromDB(limit int) []*entity.DailyCall1v1InitiateStat {
	list := make([]*entity.DailyCall1v1InitiateStat, 0, limit)
	_ = g.Model(string(entity.TbDailyCall1v1InitiateStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyCall1v1InitiateStats(list)
	return list
}

func reverseDailyCall1v1InitiateStats(list []*entity.DailyCall1v1InitiateStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyCall1v1InitiateStat(data *entity.DailyCall1v1InitiateStat) {
	if data == nil || data.ID == "" || dailyCall1v1InitiateStatCacheMgr == nil {
		return
	}
	dailyCall1v1InitiateStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
