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
	weeklyCall1v1InitiateStatCacheMgr     *cache.RowCache[*entity.WeeklyCall1v1InitiateStat]
	weeklyCall1v1InitiateStatListCacheMgr *cache.ListCache[*entity.WeeklyCall1v1InitiateStat]
)

func initWeeklyCall1v1InitiateStatDao() {
	weeklyCall1v1InitiateStatCacheMgr = cache.NewRowCache[*entity.WeeklyCall1v1InitiateStat]()
	weeklyCall1v1InitiateStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyCall1v1InitiateStat]()
}

func GetWeeklyCall1v1InitiateStatByWeek(week string) *entity.WeeklyCall1v1InitiateStat {
	return weeklyCall1v1InitiateStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyCall1v1InitiateStat, error) {
		var data *entity.WeeklyCall1v1InitiateStat
		_ = g.Model(string(entity.TbWeeklyCall1v1InitiateStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyCall1v1InitiateStat(week), nil
	})
}

func ListRecentWeeklyCall1v1InitiateStats(limit int) []*entity.WeeklyCall1v1InitiateStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyCall1v1InitiateStatKey(time.Now()))
	return weeklyCall1v1InitiateStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyCall1v1InitiateStat, error) {
		return loadRecentWeeklyCall1v1InitiateStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyCall1v1InitiateStatsFromDB(limit int) []*entity.WeeklyCall1v1InitiateStat {
	list := make([]*entity.WeeklyCall1v1InitiateStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyCall1v1InitiateStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyCall1v1InitiateStats(list)
	return list
}

func reverseWeeklyCall1v1InitiateStats(list []*entity.WeeklyCall1v1InitiateStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyCall1v1InitiateStat(data *entity.WeeklyCall1v1InitiateStat) {
	if data == nil || data.ID == "" || weeklyCall1v1InitiateStatCacheMgr == nil {
		return
	}
	weeklyCall1v1InitiateStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
