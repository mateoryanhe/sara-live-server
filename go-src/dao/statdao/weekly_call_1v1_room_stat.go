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
	weeklyCall1v1RoomStatCacheMgr     *cache.RowCache[*entity.WeeklyCall1v1RoomStat]
	weeklyCall1v1RoomStatListCacheMgr *cache.ListCache[*entity.WeeklyCall1v1RoomStat]
)

func initWeeklyCall1v1RoomStatDao() {
	weeklyCall1v1RoomStatCacheMgr = cache.NewRowCache[*entity.WeeklyCall1v1RoomStat]()
	weeklyCall1v1RoomStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyCall1v1RoomStat]()
}

func GetWeeklyCall1v1RoomStatByWeek(week string) *entity.WeeklyCall1v1RoomStat {
	return weeklyCall1v1RoomStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyCall1v1RoomStat, error) {
		var data *entity.WeeklyCall1v1RoomStat
		_ = g.Model(string(entity.TbWeeklyCall1v1RoomStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyCall1v1RoomStat(week), nil
	})
}

func ListRecentWeeklyCall1v1RoomStats(limit int) []*entity.WeeklyCall1v1RoomStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyCall1v1RoomStatKey(time.Now()))
	return weeklyCall1v1RoomStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyCall1v1RoomStat, error) {
		return loadRecentWeeklyCall1v1RoomStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyCall1v1RoomStatsFromDB(limit int) []*entity.WeeklyCall1v1RoomStat {
	list := make([]*entity.WeeklyCall1v1RoomStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyCall1v1RoomStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyCall1v1RoomStats(list)
	return list
}

func reverseWeeklyCall1v1RoomStats(list []*entity.WeeklyCall1v1RoomStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyCall1v1RoomStat(data *entity.WeeklyCall1v1RoomStat) {
	if data == nil || data.ID == "" || weeklyCall1v1RoomStatCacheMgr == nil {
		return
	}
	weeklyCall1v1RoomStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
