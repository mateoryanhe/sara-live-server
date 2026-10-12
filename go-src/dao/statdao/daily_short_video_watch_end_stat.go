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
	dailyShortVideoWatchEndStatCacheMgr     *cache.RowCache[*entity.DailyShortVideoWatchEndStat]
	dailyShortVideoWatchEndStatListCacheMgr *cache.ListCache[*entity.DailyShortVideoWatchEndStat]
)

func initDailyShortVideoWatchEndStatDao() {
	dailyShortVideoWatchEndStatCacheMgr = cache.NewRowCache[*entity.DailyShortVideoWatchEndStat]()
	dailyShortVideoWatchEndStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyShortVideoWatchEndStat]()
}

func GetDailyShortVideoWatchEndStatByDate(date string) *entity.DailyShortVideoWatchEndStat {
	return dailyShortVideoWatchEndStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyShortVideoWatchEndStat, error) {
		var data *entity.DailyShortVideoWatchEndStat
		_ = g.Model(string(entity.TbDailyShortVideoWatchEndStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyShortVideoWatchEndStat(date), nil
	})
}

func ListRecentDailyShortVideoWatchEndStats(limit int) []*entity.DailyShortVideoWatchEndStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyShortVideoWatchEndStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyShortVideoWatchEndStat, error) {
		return loadRecentDailyShortVideoWatchEndStatsFromDB(limit), nil
	})
}

func loadRecentDailyShortVideoWatchEndStatsFromDB(limit int) []*entity.DailyShortVideoWatchEndStat {
	list := make([]*entity.DailyShortVideoWatchEndStat, 0, limit)
	_ = g.Model(string(entity.TbDailyShortVideoWatchEndStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyShortVideoWatchEndStats(list)
	return list
}

func reverseDailyShortVideoWatchEndStats(list []*entity.DailyShortVideoWatchEndStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyShortVideoWatchEndStat(data *entity.DailyShortVideoWatchEndStat) {
	if data == nil || data.ID == "" || dailyShortVideoWatchEndStatCacheMgr == nil {
		return
	}
	dailyShortVideoWatchEndStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
