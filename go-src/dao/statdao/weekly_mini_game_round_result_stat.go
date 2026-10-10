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
	weeklyMiniGameRoundResultStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameRoundResultStat]
	weeklyMiniGameRoundResultStatListCacheMgr *cache.ListCache[*entity.WeeklyMiniGameRoundResultStat]
)

func initWeeklyMiniGameRoundResultStatDao() {
	weeklyMiniGameRoundResultStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameRoundResultStat]()
	weeklyMiniGameRoundResultStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyMiniGameRoundResultStat]()
}

func GetWeeklyMiniGameRoundResultStatByWeek(week string) *entity.WeeklyMiniGameRoundResultStat {
	return weeklyMiniGameRoundResultStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameRoundResultStat, error) {
		var data *entity.WeeklyMiniGameRoundResultStat
		_ = g.Model(string(entity.TbWeeklyMiniGameRoundResultStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameRoundResultStat(week), nil
	})
}

func ListRecentWeeklyMiniGameRoundResultStats(limit int) []*entity.WeeklyMiniGameRoundResultStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyMiniGameRoundResultStatKey(time.Now()))
	return weeklyMiniGameRoundResultStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyMiniGameRoundResultStat, error) {
		return loadRecentWeeklyMiniGameRoundResultStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyMiniGameRoundResultStatsFromDB(limit int) []*entity.WeeklyMiniGameRoundResultStat {
	list := make([]*entity.WeeklyMiniGameRoundResultStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameRoundResultStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameRoundResultStats(list)
	return list
}

func reverseWeeklyMiniGameRoundResultStats(list []*entity.WeeklyMiniGameRoundResultStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameRoundResultStat(data *entity.WeeklyMiniGameRoundResultStat) {
	if data == nil || data.ID == "" || weeklyMiniGameRoundResultStatCacheMgr == nil {
		return
	}
	weeklyMiniGameRoundResultStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
