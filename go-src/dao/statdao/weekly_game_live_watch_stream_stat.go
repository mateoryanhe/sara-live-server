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
	weeklyGameLiveWatchStreamStatCacheMgr     *cache.RowCache[*entity.WeeklyGameLiveWatchStreamStat]
	weeklyGameLiveWatchStreamStatListCacheMgr *cache.ListCache[*entity.WeeklyGameLiveWatchStreamStat]
)

func initWeeklyGameLiveWatchStreamStatDao() {
	weeklyGameLiveWatchStreamStatCacheMgr = cache.NewRowCache[*entity.WeeklyGameLiveWatchStreamStat]()
	weeklyGameLiveWatchStreamStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyGameLiveWatchStreamStat]()
}

func GetWeeklyGameLiveWatchStreamStatByWeek(week string) *entity.WeeklyGameLiveWatchStreamStat {
	return weeklyGameLiveWatchStreamStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyGameLiveWatchStreamStat, error) {
		var data *entity.WeeklyGameLiveWatchStreamStat
		_ = g.Model(string(entity.TbWeeklyGameLiveWatchStreamStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyGameLiveWatchStreamStat(week), nil
	})
}

func ListRecentWeeklyGameLiveWatchStreamStats(limit int) []*entity.WeeklyGameLiveWatchStreamStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyGameLiveWatchStreamStatKey(time.Now()))
	return weeklyGameLiveWatchStreamStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyGameLiveWatchStreamStat, error) {
		return loadRecentWeeklyGameLiveWatchStreamStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyGameLiveWatchStreamStatsFromDB(limit int) []*entity.WeeklyGameLiveWatchStreamStat {
	list := make([]*entity.WeeklyGameLiveWatchStreamStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyGameLiveWatchStreamStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyGameLiveWatchStreamStats(list)
	return list
}

func reverseWeeklyGameLiveWatchStreamStats(list []*entity.WeeklyGameLiveWatchStreamStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyGameLiveWatchStreamStat(data *entity.WeeklyGameLiveWatchStreamStat) {
	if data == nil || data.ID == "" || weeklyGameLiveWatchStreamStatCacheMgr == nil {
		return
	}
	weeklyGameLiveWatchStreamStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
