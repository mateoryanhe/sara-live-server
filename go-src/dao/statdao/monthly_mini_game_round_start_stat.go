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
	monthlyMiniGameRoundStartStatCacheMgr     *cache.RowCache[*entity.MonthlyMiniGameRoundStartStat]
	monthlyMiniGameRoundStartStatListCacheMgr *cache.ListCache[*entity.MonthlyMiniGameRoundStartStat]
)

func initMonthlyMiniGameRoundStartStatDao() {
	monthlyMiniGameRoundStartStatCacheMgr = cache.NewRowCache[*entity.MonthlyMiniGameRoundStartStat]()
	monthlyMiniGameRoundStartStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyMiniGameRoundStartStat]()
}

func GetMonthlyMiniGameRoundStartStatByMonth(month string) *entity.MonthlyMiniGameRoundStartStat {
	return monthlyMiniGameRoundStartStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyMiniGameRoundStartStat, error) {
		var data *entity.MonthlyMiniGameRoundStartStat
		_ = g.Model(string(entity.TbMonthlyMiniGameRoundStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyMiniGameRoundStartStat(month), nil
	})
}

func ListRecentMonthlyMiniGameRoundStartStats(limit int) []*entity.MonthlyMiniGameRoundStartStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyMiniGameRoundStartStatKey(time.Now()))
	return monthlyMiniGameRoundStartStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyMiniGameRoundStartStat, error) {
		return loadRecentMonthlyMiniGameRoundStartStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyMiniGameRoundStartStatsFromDB(limit int) []*entity.MonthlyMiniGameRoundStartStat {
	list := make([]*entity.MonthlyMiniGameRoundStartStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyMiniGameRoundStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyMiniGameRoundStartStats(list)
	return list
}

func reverseMonthlyMiniGameRoundStartStats(list []*entity.MonthlyMiniGameRoundStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyMiniGameRoundStartStat(data *entity.MonthlyMiniGameRoundStartStat) {
	if data == nil || data.ID == "" || monthlyMiniGameRoundStartStatCacheMgr == nil {
		return
	}
	monthlyMiniGameRoundStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
