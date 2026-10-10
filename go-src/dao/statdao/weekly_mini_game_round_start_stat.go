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
	weeklyMiniGameRoundStartStatCacheMgr     *cache.RowCache[*entity.WeeklyMiniGameRoundStartStat]
	weeklyMiniGameRoundStartStatListCacheMgr *cache.ListCache[*entity.WeeklyMiniGameRoundStartStat]
)

func initWeeklyMiniGameRoundStartStatDao() {
	weeklyMiniGameRoundStartStatCacheMgr = cache.NewRowCache[*entity.WeeklyMiniGameRoundStartStat]()
	weeklyMiniGameRoundStartStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyMiniGameRoundStartStat]()
}

func GetWeeklyMiniGameRoundStartStatByWeek(week string) *entity.WeeklyMiniGameRoundStartStat {
	return weeklyMiniGameRoundStartStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyMiniGameRoundStartStat, error) {
		var data *entity.WeeklyMiniGameRoundStartStat
		_ = g.Model(string(entity.TbWeeklyMiniGameRoundStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyMiniGameRoundStartStat(week), nil
	})
}

func ListRecentWeeklyMiniGameRoundStartStats(limit int) []*entity.WeeklyMiniGameRoundStartStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyMiniGameRoundStartStatKey(time.Now()))
	return weeklyMiniGameRoundStartStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyMiniGameRoundStartStat, error) {
		return loadRecentWeeklyMiniGameRoundStartStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyMiniGameRoundStartStatsFromDB(limit int) []*entity.WeeklyMiniGameRoundStartStat {
	list := make([]*entity.WeeklyMiniGameRoundStartStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyMiniGameRoundStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyMiniGameRoundStartStats(list)
	return list
}

func reverseWeeklyMiniGameRoundStartStats(list []*entity.WeeklyMiniGameRoundStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyMiniGameRoundStartStat(data *entity.WeeklyMiniGameRoundStartStat) {
	if data == nil || data.ID == "" || weeklyMiniGameRoundStartStatCacheMgr == nil {
		return
	}
	weeklyMiniGameRoundStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
