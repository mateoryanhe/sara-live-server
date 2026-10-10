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
	monthlyCall1v1InitiateStatCacheMgr     *cache.RowCache[*entity.MonthlyCall1v1InitiateStat]
	monthlyCall1v1InitiateStatListCacheMgr *cache.ListCache[*entity.MonthlyCall1v1InitiateStat]
)

func initMonthlyCall1v1InitiateStatDao() {
	monthlyCall1v1InitiateStatCacheMgr = cache.NewRowCache[*entity.MonthlyCall1v1InitiateStat]()
	monthlyCall1v1InitiateStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyCall1v1InitiateStat]()
}

func GetMonthlyCall1v1InitiateStatByMonth(month string) *entity.MonthlyCall1v1InitiateStat {
	return monthlyCall1v1InitiateStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyCall1v1InitiateStat, error) {
		var data *entity.MonthlyCall1v1InitiateStat
		_ = g.Model(string(entity.TbMonthlyCall1v1InitiateStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyCall1v1InitiateStat(month), nil
	})
}

func ListRecentMonthlyCall1v1InitiateStats(limit int) []*entity.MonthlyCall1v1InitiateStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyCall1v1InitiateStatKey(time.Now()))
	return monthlyCall1v1InitiateStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyCall1v1InitiateStat, error) {
		return loadRecentMonthlyCall1v1InitiateStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyCall1v1InitiateStatsFromDB(limit int) []*entity.MonthlyCall1v1InitiateStat {
	list := make([]*entity.MonthlyCall1v1InitiateStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyCall1v1InitiateStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyCall1v1InitiateStats(list)
	return list
}

func reverseMonthlyCall1v1InitiateStats(list []*entity.MonthlyCall1v1InitiateStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyCall1v1InitiateStat(data *entity.MonthlyCall1v1InitiateStat) {
	if data == nil || data.ID == "" || monthlyCall1v1InitiateStatCacheMgr == nil {
		return
	}
	monthlyCall1v1InitiateStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
