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
	monthlyMiniGameExposureStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameExposureStat]
	monthlyMiniGameExposureStatListCacheMgr *cache.ListCache[*entity.MonthlyMiniGameExposureStat]
)

func initMonthlyMiniGameExposureStatDao() {
	monthlyMiniGameExposureStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameExposureStat]()
	monthlyMiniGameExposureStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyMiniGameExposureStat]()
}

func GetMonthlyMiniGameExposureStatByMonth(month string) *entity.MonthlyMiniGameExposureStat {
	return monthlyMiniGameExposureStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameExposureStat, error) {
		var data *entity.MonthlyMiniGameExposureStat
		_ = g.Model(string(entity.TbMonthlyMiniGameExposureStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameExposureStat(month), nil
	})
}

func ListRecentMonthlyMiniGameExposureStats(limit int) []*entity.MonthlyMiniGameExposureStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyMiniGameExposureStatKey(time.Now()))
	return monthlyMiniGameExposureStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyMiniGameExposureStat, error) {
		return loadRecentMonthlyMiniGameExposureStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyMiniGameExposureStatsFromDB(limit int) []*entity.MonthlyMiniGameExposureStat {
	list := make([]*entity.MonthlyMiniGameExposureStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameExposureStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameExposureStats(list)
	return list
}

func reverseMonthlyMiniGameExposureStats(list []*entity.MonthlyMiniGameExposureStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameExposureStat(data *entity.MonthlyMiniGameExposureStat) {
	if data == nil || data.ID == "" || monthlyMiniGameExposureStatCacheMgr == nil {
		return
	}
	monthlyMiniGameExposureStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
