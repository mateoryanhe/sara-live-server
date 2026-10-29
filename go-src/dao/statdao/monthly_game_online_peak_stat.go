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
	monthlyGameOnlinePeakStatCacheMgr     *cache.RowCache[*entity.MonthlyGameOnlinePeakStat]
	monthlyGameOnlinePeakStatListCacheMgr *cache.ListCache[*entity.MonthlyGameOnlinePeakStat]
)

func initMonthlyGameOnlinePeakStatDao() {
	monthlyGameOnlinePeakStatCacheMgr = cache.NewRowCache[*entity.MonthlyGameOnlinePeakStat]()
	monthlyGameOnlinePeakStatListCacheMgr = cache.NewDashboardListCache[*entity.MonthlyGameOnlinePeakStat]()
}

func GetMonthlyGameOnlinePeakStatByMonth(month string) *entity.MonthlyGameOnlinePeakStat {
	return monthlyGameOnlinePeakStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyGameOnlinePeakStat, error) {
		var data *entity.MonthlyGameOnlinePeakStat
		_ = g.Model(string(entity.TbMonthlyGameOnlinePeakStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyGameOnlinePeakStat(month), nil
	})
}

func ListRecentMonthlyGameOnlinePeakStats(limit int) []*entity.MonthlyGameOnlinePeakStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyGameBetUserStatKey(time.Now()))
	return monthlyGameOnlinePeakStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyGameOnlinePeakStat, error) {
		return loadRecentMonthlyGameOnlinePeakStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyGameOnlinePeakStatsFromDB(limit int) []*entity.MonthlyGameOnlinePeakStat {
	list := make([]*entity.MonthlyGameOnlinePeakStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyGameOnlinePeakStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyGameOnlinePeakStats(list)
	return list
}

func reverseMonthlyGameOnlinePeakStats(list []*entity.MonthlyGameOnlinePeakStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyGameOnlinePeakStat(data *entity.MonthlyGameOnlinePeakStat) {
	if data == nil || data.ID == "" || monthlyGameOnlinePeakStatCacheMgr == nil {
		return
	}
	monthlyGameOnlinePeakStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
