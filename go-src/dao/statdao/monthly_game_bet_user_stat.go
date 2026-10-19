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
	monthlyGameBetUserStatCacheMgr     *cache.RowCache[*entity.MonthlyGameBetUserStat]
	monthlyGameBetUserStatListCacheMgr *cache.ListCache[*entity.MonthlyGameBetUserStat]
)

func initMonthlyGameBetUserStatDao() {
	monthlyGameBetUserStatCacheMgr = cache.NewRowCache[*entity.MonthlyGameBetUserStat]()
	monthlyGameBetUserStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyGameBetUserStat]()
}

func GetMonthlyGameBetUserStatByMonth(month string) *entity.MonthlyGameBetUserStat {
	return monthlyGameBetUserStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyGameBetUserStat, error) {
		var data *entity.MonthlyGameBetUserStat
		_ = g.Model(string(entity.TbMonthlyGameBetUserStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyGameBetUserStat(month), nil
	})
}

func ListRecentMonthlyGameBetUserStats(limit int) []*entity.MonthlyGameBetUserStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyGameBetUserStatKey(time.Now()))
	return monthlyGameBetUserStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyGameBetUserStat, error) {
		return loadRecentMonthlyGameBetUserStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyGameBetUserStatsFromDB(limit int) []*entity.MonthlyGameBetUserStat {
	list := make([]*entity.MonthlyGameBetUserStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyGameBetUserStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyGameBetUserStats(list)
	return list
}

func reverseMonthlyGameBetUserStats(list []*entity.MonthlyGameBetUserStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyGameBetUserStat(data *entity.MonthlyGameBetUserStat) {
	if data == nil || data.ID == "" || monthlyGameBetUserStatCacheMgr == nil {
		return
	}
	monthlyGameBetUserStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
