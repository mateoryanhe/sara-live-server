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
	dailyShortVideoPaySuccessStatCacheMgr     *cache.RowCache[*entity.DailyShortVideoPaySuccessStat]
)

func initDailyShortVideoPaySuccessStatDao() {
	dailyShortVideoPaySuccessStatCacheMgr = cache.NewRowCache[*entity.DailyShortVideoPaySuccessStat]()
}

func GetDailyShortVideoPaySuccessStatByDate(date string) *entity.DailyShortVideoPaySuccessStat {
	return dailyShortVideoPaySuccessStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyShortVideoPaySuccessStat, error) {
		var data *entity.DailyShortVideoPaySuccessStat
		_ = g.Model(string(entity.TbDailyShortVideoPaySuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyShortVideoPaySuccessStat(date), nil
	})
}

func ListRecentDailyShortVideoPaySuccessStats(limit int) []*entity.DailyShortVideoPaySuccessStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyShortVideoPaySuccessStatsFromDB(limit)
}

func loadRecentDailyShortVideoPaySuccessStatsFromDB(limit int) []*entity.DailyShortVideoPaySuccessStat {
	list := make([]*entity.DailyShortVideoPaySuccessStat, 0, limit)
	_ = g.Model(string(entity.TbDailyShortVideoPaySuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyShortVideoPaySuccessStats(list)
	return list
}

func reverseDailyShortVideoPaySuccessStats(list []*entity.DailyShortVideoPaySuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyShortVideoPaySuccessStat(data *entity.DailyShortVideoPaySuccessStat) {
	if data == nil || data.ID == "" || dailyShortVideoPaySuccessStatCacheMgr == nil {
		return
	}
	dailyShortVideoPaySuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
