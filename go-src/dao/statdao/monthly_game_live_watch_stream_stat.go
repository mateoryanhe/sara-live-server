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
	monthlyGameLiveWatchStreamStatCacheMgr     *cache.RowCache[*entity.MonthlyGameLiveWatchStreamStat]
	monthlyGameLiveWatchStreamStatListCacheMgr *cache.ListCache[*entity.MonthlyGameLiveWatchStreamStat]
)

func initMonthlyGameLiveWatchStreamStatDao() {
	monthlyGameLiveWatchStreamStatCacheMgr = cache.NewRowCache[*entity.MonthlyGameLiveWatchStreamStat]()
	monthlyGameLiveWatchStreamStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyGameLiveWatchStreamStat]()
}

func GetMonthlyGameLiveWatchStreamStatByMonth(month string) *entity.MonthlyGameLiveWatchStreamStat {
	return monthlyGameLiveWatchStreamStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyGameLiveWatchStreamStat, error) {
		var data *entity.MonthlyGameLiveWatchStreamStat
		_ = g.Model(string(entity.TbMonthlyGameLiveWatchStreamStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyGameLiveWatchStreamStat(month), nil
	})
}

func ListRecentMonthlyGameLiveWatchStreamStats(limit int) []*entity.MonthlyGameLiveWatchStreamStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyGameLiveWatchStreamStatKey(time.Now()))
	return monthlyGameLiveWatchStreamStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyGameLiveWatchStreamStat, error) {
		return loadRecentMonthlyGameLiveWatchStreamStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyGameLiveWatchStreamStatsFromDB(limit int) []*entity.MonthlyGameLiveWatchStreamStat {
	list := make([]*entity.MonthlyGameLiveWatchStreamStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyGameLiveWatchStreamStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyGameLiveWatchStreamStats(list)
	return list
}

func reverseMonthlyGameLiveWatchStreamStats(list []*entity.MonthlyGameLiveWatchStreamStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyGameLiveWatchStreamStat(data *entity.MonthlyGameLiveWatchStreamStat) {
	if data == nil || data.ID == "" || monthlyGameLiveWatchStreamStatCacheMgr == nil {
		return
	}
	monthlyGameLiveWatchStreamStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
