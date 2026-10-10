package statdao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var (
	monthlyMiniGameWebViewLoadStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameWebViewLoadStat]
	monthlyMiniGameWebViewLoadStatListCacheMgr *cache.ListCache[*entity.MonthlyMiniGameWebViewLoadStat]
)

func initMonthlyMiniGameWebViewLoadStatDao() {
	monthlyMiniGameWebViewLoadStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameWebViewLoadStat]()
	monthlyMiniGameWebViewLoadStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyMiniGameWebViewLoadStat]()
}

func GetMonthlyMiniGameWebViewLoadStatByMonth(month string) *entity.MonthlyMiniGameWebViewLoadStat {
	return monthlyMiniGameWebViewLoadStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameWebViewLoadStat, error) {
		var data *entity.MonthlyMiniGameWebViewLoadStat
		_ = g.Model(string(entity.TbMonthlyMiniGameWebViewLoadStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameWebViewLoadStat(month), nil
	})
}

func ListRecentMonthlyMiniGameWebViewLoadStats(limit int) []*entity.MonthlyMiniGameWebViewLoadStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyMiniGameWebViewLoadStatKey(time.Now()))
	return monthlyMiniGameWebViewLoadStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyMiniGameWebViewLoadStat, error) {
		return loadRecentMonthlyMiniGameWebViewLoadStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyMiniGameWebViewLoadStatsFromDB(limit int) []*entity.MonthlyMiniGameWebViewLoadStat {
	list := make([]*entity.MonthlyMiniGameWebViewLoadStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameWebViewLoadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameWebViewLoadStats(list)
	return list
}

func reverseMonthlyMiniGameWebViewLoadStats(list []*entity.MonthlyMiniGameWebViewLoadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameWebViewLoadStat(data *entity.MonthlyMiniGameWebViewLoadStat) {
	if data == nil || data.ID == "" || monthlyMiniGameWebViewLoadStatCacheMgr == nil {
		return
	}
	monthlyMiniGameWebViewLoadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
