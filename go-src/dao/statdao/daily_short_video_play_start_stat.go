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
	dailyShortVideoPlayStartStatCacheMgr     *cache.RowCache[*entity.DailyShortVideoPlayStartStat]
)

func initDailyShortVideoPlayStartStatDao() {
	dailyShortVideoPlayStartStatCacheMgr = cache.NewRowCache[*entity.DailyShortVideoPlayStartStat]()
}

func GetDailyShortVideoPlayStartStatByDate(date string) *entity.DailyShortVideoPlayStartStat {
	return dailyShortVideoPlayStartStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyShortVideoPlayStartStat, error) {
		var data *entity.DailyShortVideoPlayStartStat
		_ = g.Model(string(entity.TbDailyShortVideoPlayStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyShortVideoPlayStartStat(date), nil
	})
}

func ListRecentDailyShortVideoPlayStartStats(limit int) []*entity.DailyShortVideoPlayStartStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyShortVideoPlayStartStatsFromDB(limit)
}

func loadRecentDailyShortVideoPlayStartStatsFromDB(limit int) []*entity.DailyShortVideoPlayStartStat {
	list := make([]*entity.DailyShortVideoPlayStartStat, 0, limit)
	_ = g.Model(string(entity.TbDailyShortVideoPlayStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyShortVideoPlayStartStats(list)
	return list
}

func reverseDailyShortVideoPlayStartStats(list []*entity.DailyShortVideoPlayStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyShortVideoPlayStartStat(data *entity.DailyShortVideoPlayStartStat) {
	if data == nil || data.ID == "" || dailyShortVideoPlayStartStatCacheMgr == nil {
		return
	}
	dailyShortVideoPlayStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
