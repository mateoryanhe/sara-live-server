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
	weeklyH5DownloadStatCacheMgr     *cache.RowCache[*entity.WeeklyH5DownloadStat]
	weeklyH5DownloadStatListCacheMgr *cache.ListCache[*entity.WeeklyH5DownloadStat]
)

func initWeeklyH5DownloadStatDao() {
	weeklyH5DownloadStatCacheMgr = cache.NewRowCache[*entity.WeeklyH5DownloadStat]()
	weeklyH5DownloadStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyH5DownloadStat]()
}

func GetWeeklyH5DownloadStatByWeek(week string) *entity.WeeklyH5DownloadStat {
	return weeklyH5DownloadStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyH5DownloadStat, error) {
		var data *entity.WeeklyH5DownloadStat
		_ = g.Model(string(entity.TbWeeklyH5DownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyH5DownloadStat(week), nil
	})
}

func ListRecentWeeklyH5DownloadStats(limit int) []*entity.WeeklyH5DownloadStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyH5DownloadStatKey(time.Now()))
	return weeklyH5DownloadStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyH5DownloadStat, error) {
		return loadRecentWeeklyH5DownloadStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyH5DownloadStatsFromDB(limit int) []*entity.WeeklyH5DownloadStat {
	list := make([]*entity.WeeklyH5DownloadStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyH5DownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyH5DownloadStats(list)
	return list
}

func reverseWeeklyH5DownloadStats(list []*entity.WeeklyH5DownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyH5DownloadStat(data *entity.WeeklyH5DownloadStat) {
	if data == nil || data.ID == "" || weeklyH5DownloadStatCacheMgr == nil {
		return
	}
	weeklyH5DownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
