package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var (
	weeklyOneToOneRoomListStatCacheMgr     *cache.RowCache[*entity.WeeklyOneToOneRoomListStat]
)

func initWeeklyOneToOneRoomListStatDao() {
	weeklyOneToOneRoomListStatCacheMgr = cache.NewRowCache[*entity.WeeklyOneToOneRoomListStat]()
}

func GetWeeklyOneToOneRoomListStatByWeek(week string) *entity.WeeklyOneToOneRoomListStat {
	return weeklyOneToOneRoomListStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyOneToOneRoomListStat, error) {
		var data *entity.WeeklyOneToOneRoomListStat
		_ = g.Model(string(entity.TbWeeklyOneToOneRoomListStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyOneToOneRoomListStat(week), nil
	})
}

func ListRecentWeeklyOneToOneRoomListStats(limit int) []*entity.WeeklyOneToOneRoomListStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyOneToOneRoomListStatsFromDB(limit)
}

func loadRecentWeeklyOneToOneRoomListStatsFromDB(limit int) []*entity.WeeklyOneToOneRoomListStat {
	list := make([]*entity.WeeklyOneToOneRoomListStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyOneToOneRoomListStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyOneToOneRoomListStats(list)
	return list
}

func reverseWeeklyOneToOneRoomListStats(list []*entity.WeeklyOneToOneRoomListStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyOneToOneRoomListStat(data *entity.WeeklyOneToOneRoomListStat) {
	if data == nil || data.ID == "" || weeklyOneToOneRoomListStatCacheMgr == nil {
		return
	}
	weeklyOneToOneRoomListStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
