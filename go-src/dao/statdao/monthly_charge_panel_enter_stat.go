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
	monthlyChargePanelEnterStatCacheMgr     *cache.RowCache[*entity.MonthlyChargePanelEnterStat]
	monthlyChargePanelEnterStatListCacheMgr *cache.ListCache[*entity.MonthlyChargePanelEnterStat]
)

func initMonthlyChargePanelEnterStatDao() {
	monthlyChargePanelEnterStatCacheMgr = cache.NewRowCache[*entity.MonthlyChargePanelEnterStat]()
	monthlyChargePanelEnterStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyChargePanelEnterStat]()
}

func GetMonthlyChargePanelEnterStatByMonth(month string) *entity.MonthlyChargePanelEnterStat {
	return monthlyChargePanelEnterStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyChargePanelEnterStat, error) {
		var data *entity.MonthlyChargePanelEnterStat
		_ = g.Model(string(entity.TbMonthlyChargePanelEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyChargePanelEnterStat(month), nil
	})
}

func ListRecentMonthlyChargePanelEnterStats(limit int) []*entity.MonthlyChargePanelEnterStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyChargePanelEnterStatKey(time.Now()))
	return monthlyChargePanelEnterStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyChargePanelEnterStat, error) {
		return loadRecentMonthlyChargePanelEnterStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyChargePanelEnterStatsFromDB(limit int) []*entity.MonthlyChargePanelEnterStat {
	list := make([]*entity.MonthlyChargePanelEnterStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyChargePanelEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyChargePanelEnterStats(list)
	return list
}

func reverseMonthlyChargePanelEnterStats(list []*entity.MonthlyChargePanelEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyChargePanelEnterStat(data *entity.MonthlyChargePanelEnterStat) {
	if data == nil || data.ID == "" || monthlyChargePanelEnterStatCacheMgr == nil {
		return
	}
	monthlyChargePanelEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
