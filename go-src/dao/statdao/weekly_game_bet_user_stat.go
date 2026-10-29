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
	weeklyGameBetUserStatCacheMgr     *cache.RowCache[*entity.WeeklyGameBetUserStat]
	weeklyGameBetUserStatListCacheMgr *cache.ListCache[*entity.WeeklyGameBetUserStat]
)

func initWeeklyGameBetUserStatDao() {
	weeklyGameBetUserStatCacheMgr = cache.NewRowCache[*entity.WeeklyGameBetUserStat]()
	weeklyGameBetUserStatListCacheMgr = cache.NewDashboardListCache[*entity.WeeklyGameBetUserStat]()
}

func GetWeeklyGameBetUserStatByWeek(week string) *entity.WeeklyGameBetUserStat {
	return weeklyGameBetUserStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyGameBetUserStat, error) {
		var data *entity.WeeklyGameBetUserStat
		_ = g.Model(string(entity.TbWeeklyGameBetUserStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyGameBetUserStat(week), nil
	})
}

func ListRecentWeeklyGameBetUserStats(limit int) []*entity.WeeklyGameBetUserStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyLoginStatKey(time.Now()))
	return weeklyGameBetUserStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyGameBetUserStat, error) {
		return loadRecentWeeklyGameBetUserStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyGameBetUserStatsFromDB(limit int) []*entity.WeeklyGameBetUserStat {
	list := make([]*entity.WeeklyGameBetUserStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyGameBetUserStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyGameBetUserStats(list)
	return list
}

func reverseWeeklyGameBetUserStats(list []*entity.WeeklyGameBetUserStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyGameBetUserStat(data *entity.WeeklyGameBetUserStat) {
	if data == nil || data.ID == "" || weeklyGameBetUserStatCacheMgr == nil {
		return
	}
	weeklyGameBetUserStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
