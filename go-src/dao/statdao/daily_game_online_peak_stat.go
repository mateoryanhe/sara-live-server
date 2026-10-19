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
	dailyGameOnlinePeakStatCacheMgr     *cache.RowCache[*entity.DailyGameOnlinePeakStat]
	dailyGameOnlinePeakStatListCacheMgr *cache.ListCache[*entity.DailyGameOnlinePeakStat]
)

func initDailyGameOnlinePeakStatDao() {
	dailyGameOnlinePeakStatCacheMgr = cache.NewRowCache[*entity.DailyGameOnlinePeakStat]()
	dailyGameOnlinePeakStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyGameOnlinePeakStat]()
}

func GetDailyGameOnlinePeakStatByDate(date string) *entity.DailyGameOnlinePeakStat {
	return dailyGameOnlinePeakStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGameOnlinePeakStat, error) {
		var data *entity.DailyGameOnlinePeakStat
		_ = g.Model(string(entity.TbDailyGameOnlinePeakStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGameOnlinePeakStat(date), nil
	})
}

func ListRecentDailyGameOnlinePeakStats(limit int) []*entity.DailyGameOnlinePeakStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyGameOnlinePeakStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyGameOnlinePeakStat, error) {
		return loadRecentDailyGameOnlinePeakStatsFromDB(limit), nil
	})
}

func loadRecentDailyGameOnlinePeakStatsFromDB(limit int) []*entity.DailyGameOnlinePeakStat {
	list := make([]*entity.DailyGameOnlinePeakStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGameOnlinePeakStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGameOnlinePeakStats(list)
	return list
}

func reverseDailyGameOnlinePeakStats(list []*entity.DailyGameOnlinePeakStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGameOnlinePeakStat(data *entity.DailyGameOnlinePeakStat) {
	if data == nil || data.ID == "" || dailyGameOnlinePeakStatCacheMgr == nil {
		return
	}
	dailyGameOnlinePeakStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
