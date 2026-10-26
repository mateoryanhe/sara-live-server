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
	weeklyOfficialSiteDownloadStatCacheMgr     *cache.RowCache[*entity.WeeklyOfficialSiteDownloadStat]
	weeklyOfficialSiteDownloadStatListCacheMgr *cache.ListCache[*entity.WeeklyOfficialSiteDownloadStat]
)

func initWeeklyOfficialSiteDownloadStatDao() {
	weeklyOfficialSiteDownloadStatCacheMgr = cache.NewRowCache[*entity.WeeklyOfficialSiteDownloadStat]()
	weeklyOfficialSiteDownloadStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyOfficialSiteDownloadStat]()
}

func GetWeeklyOfficialSiteDownloadStatByWeek(week string) *entity.WeeklyOfficialSiteDownloadStat {
	return weeklyOfficialSiteDownloadStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyOfficialSiteDownloadStat, error) {
		var data *entity.WeeklyOfficialSiteDownloadStat
		_ = g.Model(string(entity.TbWeeklyOfficialSiteDownloadStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyOfficialSiteDownloadStat(week), nil
	})
}

func ListRecentWeeklyOfficialSiteDownloadStats(limit int) []*entity.WeeklyOfficialSiteDownloadStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyOfficialSiteDownloadStatKey(time.Now()))
	return weeklyOfficialSiteDownloadStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyOfficialSiteDownloadStat, error) {
		return loadRecentWeeklyOfficialSiteDownloadStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyOfficialSiteDownloadStatsFromDB(limit int) []*entity.WeeklyOfficialSiteDownloadStat {
	list := make([]*entity.WeeklyOfficialSiteDownloadStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyOfficialSiteDownloadStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyOfficialSiteDownloadStats(list)
	return list
}

func reverseWeeklyOfficialSiteDownloadStats(list []*entity.WeeklyOfficialSiteDownloadStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyOfficialSiteDownloadStat(data *entity.WeeklyOfficialSiteDownloadStat) {
	if data == nil || data.ID == "" || weeklyOfficialSiteDownloadStatCacheMgr == nil {
		return
	}
	weeklyOfficialSiteDownloadStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
