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
	monthlyShareShortLinkClickStatCacheMgr     *cache.RowCache[*entity.MonthlyShareShortLinkClickStat]
	monthlyShareShortLinkClickStatListCacheMgr *cache.ListCache[*entity.MonthlyShareShortLinkClickStat]
)

func initMonthlyShareShortLinkClickStatDao() {
	monthlyShareShortLinkClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyShareShortLinkClickStat]()
	monthlyShareShortLinkClickStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyShareShortLinkClickStat]()
}

func GetMonthlyShareShortLinkClickStatByMonth(month string) *entity.MonthlyShareShortLinkClickStat {
	return monthlyShareShortLinkClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyShareShortLinkClickStat, error) {
		var data *entity.MonthlyShareShortLinkClickStat
		_ = g.Model(string(entity.TbMonthlyShareShortLinkClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyShareShortLinkClickStat(month), nil
	})
}

func ListRecentMonthlyShareShortLinkClickStats(limit int) []*entity.MonthlyShareShortLinkClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyShareShortLinkClickStatKey(time.Now()))
	return monthlyShareShortLinkClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyShareShortLinkClickStat, error) {
		return loadRecentMonthlyShareShortLinkClickStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyShareShortLinkClickStatsFromDB(limit int) []*entity.MonthlyShareShortLinkClickStat {
	list := make([]*entity.MonthlyShareShortLinkClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyShareShortLinkClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyShareShortLinkClickStats(list)
	return list
}

func reverseMonthlyShareShortLinkClickStats(list []*entity.MonthlyShareShortLinkClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyShareShortLinkClickStat(data *entity.MonthlyShareShortLinkClickStat) {
	if data == nil || data.ID == "" || monthlyShareShortLinkClickStatCacheMgr == nil {
		return
	}
	monthlyShareShortLinkClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
