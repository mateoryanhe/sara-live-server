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
	weeklyShortVideoPaySuccessStatCacheMgr     *cache.RowCache[*entity.WeeklyShortVideoPaySuccessStat]
	weeklyShortVideoPaySuccessStatListCacheMgr *cache.ListCache[*entity.WeeklyShortVideoPaySuccessStat]
)

func initWeeklyShortVideoPaySuccessStatDao() {
	weeklyShortVideoPaySuccessStatCacheMgr = cache.NewRowCache[*entity.WeeklyShortVideoPaySuccessStat]()
	weeklyShortVideoPaySuccessStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyShortVideoPaySuccessStat]()
}

func GetWeeklyShortVideoPaySuccessStatByWeek(week string) *entity.WeeklyShortVideoPaySuccessStat {
	return weeklyShortVideoPaySuccessStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyShortVideoPaySuccessStat, error) {
		var data *entity.WeeklyShortVideoPaySuccessStat
		_ = g.Model(string(entity.TbWeeklyShortVideoPaySuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyShortVideoPaySuccessStat(week), nil
	})
}

func ListRecentWeeklyShortVideoPaySuccessStats(limit int) []*entity.WeeklyShortVideoPaySuccessStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyShortVideoPaySuccessStatKey(time.Now()))
	return weeklyShortVideoPaySuccessStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyShortVideoPaySuccessStat, error) {
		return loadRecentWeeklyShortVideoPaySuccessStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyShortVideoPaySuccessStatsFromDB(limit int) []*entity.WeeklyShortVideoPaySuccessStat {
	list := make([]*entity.WeeklyShortVideoPaySuccessStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyShortVideoPaySuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyShortVideoPaySuccessStats(list)
	return list
}

func reverseWeeklyShortVideoPaySuccessStats(list []*entity.WeeklyShortVideoPaySuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyShortVideoPaySuccessStat(data *entity.WeeklyShortVideoPaySuccessStat) {
	if data == nil || data.ID == "" || weeklyShortVideoPaySuccessStatCacheMgr == nil {
		return
	}
	weeklyShortVideoPaySuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
