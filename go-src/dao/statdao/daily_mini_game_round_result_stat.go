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
	dailyMiniGameRoundResultStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameRoundResultStat]
)

func initDailyMiniGameRoundResultStatDao() {
	dailyMiniGameRoundResultStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameRoundResultStat]()
}

func GetDailyMiniGameRoundResultStatByDate(date string) *entity.DailyMiniGameRoundResultStat {
	return dailyMiniGameRoundResultStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameRoundResultStat, error) {
		var data *entity.DailyMiniGameRoundResultStat
		_ = g.Model(string(entity.TbDailyMiniGameRoundResultStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameRoundResultStat(date), nil
	})
}

func ListRecentDailyMiniGameRoundResultStats(limit int) []*entity.DailyMiniGameRoundResultStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyMiniGameRoundResultStatsFromDB(limit)
}

func loadRecentDailyMiniGameRoundResultStatsFromDB(limit int) []*entity.DailyMiniGameRoundResultStat {
	list := make([]*entity.DailyMiniGameRoundResultStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameRoundResultStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameRoundResultStats(list)
	return list
}

func reverseDailyMiniGameRoundResultStats(list []*entity.DailyMiniGameRoundResultStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameRoundResultStat(data *entity.DailyMiniGameRoundResultStat) {
	if data == nil || data.ID == "" || dailyMiniGameRoundResultStatCacheMgr == nil {
		return
	}
	dailyMiniGameRoundResultStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
