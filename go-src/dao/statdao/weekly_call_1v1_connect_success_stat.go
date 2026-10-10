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
	weeklyCall1v1ConnectSuccessStatCacheMgr     *cache.RowCache[*entity.WeeklyCall1v1ConnectSuccessStat]
	weeklyCall1v1ConnectSuccessStatListCacheMgr *cache.ListCache[*entity.WeeklyCall1v1ConnectSuccessStat]
)

func initWeeklyCall1v1ConnectSuccessStatDao() {
	weeklyCall1v1ConnectSuccessStatCacheMgr = cache.NewRowCache[*entity.WeeklyCall1v1ConnectSuccessStat]()
	weeklyCall1v1ConnectSuccessStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyCall1v1ConnectSuccessStat]()
}

func GetWeeklyCall1v1ConnectSuccessStatByWeek(week string) *entity.WeeklyCall1v1ConnectSuccessStat {
	return weeklyCall1v1ConnectSuccessStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyCall1v1ConnectSuccessStat, error) {
		var data *entity.WeeklyCall1v1ConnectSuccessStat
		_ = g.Model(string(entity.TbWeeklyCall1v1ConnectSuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyCall1v1ConnectSuccessStat(week), nil
	})
}

func ListRecentWeeklyCall1v1ConnectSuccessStats(limit int) []*entity.WeeklyCall1v1ConnectSuccessStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyCall1v1ConnectSuccessStatKey(time.Now()))
	return weeklyCall1v1ConnectSuccessStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyCall1v1ConnectSuccessStat, error) {
		return loadRecentWeeklyCall1v1ConnectSuccessStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyCall1v1ConnectSuccessStatsFromDB(limit int) []*entity.WeeklyCall1v1ConnectSuccessStat {
	list := make([]*entity.WeeklyCall1v1ConnectSuccessStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyCall1v1ConnectSuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyCall1v1ConnectSuccessStats(list)
	return list
}

func reverseWeeklyCall1v1ConnectSuccessStats(list []*entity.WeeklyCall1v1ConnectSuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyCall1v1ConnectSuccessStat(data *entity.WeeklyCall1v1ConnectSuccessStat) {
	if data == nil || data.ID == "" || weeklyCall1v1ConnectSuccessStatCacheMgr == nil {
		return
	}
	weeklyCall1v1ConnectSuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
