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
	monthlyShortVideoPlayStartStatCacheMgr     *cache.RowCache[*entity.MonthlyShortVideoPlayStartStat]
	monthlyShortVideoPlayStartStatListCacheMgr *cache.ListCache[*entity.MonthlyShortVideoPlayStartStat]
)

func initMonthlyShortVideoPlayStartStatDao() {
	monthlyShortVideoPlayStartStatCacheMgr = cache.NewRowCache[*entity.MonthlyShortVideoPlayStartStat]()
	monthlyShortVideoPlayStartStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyShortVideoPlayStartStat]()
}

func GetMonthlyShortVideoPlayStartStatByMonth(month string) *entity.MonthlyShortVideoPlayStartStat {
	return monthlyShortVideoPlayStartStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyShortVideoPlayStartStat, error) {
		var data *entity.MonthlyShortVideoPlayStartStat
		_ = g.Model(string(entity.TbMonthlyShortVideoPlayStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyShortVideoPlayStartStat(month), nil
	})
}

func ListRecentMonthlyShortVideoPlayStartStats(limit int) []*entity.MonthlyShortVideoPlayStartStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyShortVideoPlayStartStatKey(time.Now()))
	return monthlyShortVideoPlayStartStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyShortVideoPlayStartStat, error) {
		return loadRecentMonthlyShortVideoPlayStartStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyShortVideoPlayStartStatsFromDB(limit int) []*entity.MonthlyShortVideoPlayStartStat {
	list := make([]*entity.MonthlyShortVideoPlayStartStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyShortVideoPlayStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyShortVideoPlayStartStats(list)
	return list
}

func reverseMonthlyShortVideoPlayStartStats(list []*entity.MonthlyShortVideoPlayStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyShortVideoPlayStartStat(data *entity.MonthlyShortVideoPlayStartStat) {
	if data == nil || data.ID == "" || monthlyShortVideoPlayStartStatCacheMgr == nil {
		return
	}
	monthlyShortVideoPlayStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
