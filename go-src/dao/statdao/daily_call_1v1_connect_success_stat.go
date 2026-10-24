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
	dailyCall1v1ConnectSuccessStatCacheMgr *cache.RowCache[*entity.DailyCall1v1ConnectSuccessStat]
)

func initDailyCall1v1ConnectSuccessStatDao() {
	dailyCall1v1ConnectSuccessStatCacheMgr = cache.NewRowCache[*entity.DailyCall1v1ConnectSuccessStat]()
}

func GetDailyCall1v1ConnectSuccessStatByDate(date string) *entity.DailyCall1v1ConnectSuccessStat {
	return dailyCall1v1ConnectSuccessStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyCall1v1ConnectSuccessStat, error) {
		var data *entity.DailyCall1v1ConnectSuccessStat
		_ = g.Model(string(entity.TbDailyCall1v1ConnectSuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyCall1v1ConnectSuccessStat(date), nil
	})
}

func ListRecentDailyCall1v1ConnectSuccessStats(limit int) []*entity.DailyCall1v1ConnectSuccessStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyCall1v1ConnectSuccessStatsFromDB(limit)
}

func loadRecentDailyCall1v1ConnectSuccessStatsFromDB(limit int) []*entity.DailyCall1v1ConnectSuccessStat {
	list := make([]*entity.DailyCall1v1ConnectSuccessStat, 0, limit)
	_ = g.Model(string(entity.TbDailyCall1v1ConnectSuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyCall1v1ConnectSuccessStats(list)
	return list
}

func reverseDailyCall1v1ConnectSuccessStats(list []*entity.DailyCall1v1ConnectSuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyCall1v1ConnectSuccessStat(data *entity.DailyCall1v1ConnectSuccessStat) {
	if data == nil || data.ID == "" || dailyCall1v1ConnectSuccessStatCacheMgr == nil {
		return
	}
	dailyCall1v1ConnectSuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
