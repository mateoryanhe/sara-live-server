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
	monthlyMiniGameRoundResultStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameRoundResultStat]
	monthlyMiniGameRoundResultStatListCacheMgr *cache.ListCache[*entity.MonthlyMiniGameRoundResultStat]
)

func initMonthlyMiniGameRoundResultStatDao() {
	monthlyMiniGameRoundResultStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameRoundResultStat]()
	monthlyMiniGameRoundResultStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyMiniGameRoundResultStat]()
}

func GetMonthlyMiniGameRoundResultStatByMonth(month string) *entity.MonthlyMiniGameRoundResultStat {
	return monthlyMiniGameRoundResultStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameRoundResultStat, error) {
		var data *entity.MonthlyMiniGameRoundResultStat
		_ = g.Model(string(entity.TbMonthlyMiniGameRoundResultStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameRoundResultStat(month), nil
	})
}

func ListRecentMonthlyMiniGameRoundResultStats(limit int) []*entity.MonthlyMiniGameRoundResultStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyMiniGameRoundResultStatKey(time.Now()))
	return monthlyMiniGameRoundResultStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyMiniGameRoundResultStat, error) {
		return loadRecentMonthlyMiniGameRoundResultStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyMiniGameRoundResultStatsFromDB(limit int) []*entity.MonthlyMiniGameRoundResultStat {
	list := make([]*entity.MonthlyMiniGameRoundResultStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameRoundResultStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameRoundResultStats(list)
	return list
}

func reverseMonthlyMiniGameRoundResultStats(list []*entity.MonthlyMiniGameRoundResultStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameRoundResultStat(data *entity.MonthlyMiniGameRoundResultStat) {
	if data == nil || data.ID == "" || monthlyMiniGameRoundResultStatCacheMgr == nil {
		return
	}
	monthlyMiniGameRoundResultStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
