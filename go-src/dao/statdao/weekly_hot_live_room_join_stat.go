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
	weeklyHotLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.WeeklyHotLiveRoomJoinStat]
)

func initWeeklyHotLiveRoomJoinStatDao() {
	weeklyHotLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.WeeklyHotLiveRoomJoinStat]()
}

func GetWeeklyHotLiveRoomJoinStatByWeek(week string) *entity.WeeklyHotLiveRoomJoinStat {
	return weeklyHotLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyHotLiveRoomJoinStat, error) {
		var data *entity.WeeklyHotLiveRoomJoinStat
		_ = g.Model(string(entity.TbWeeklyHotLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyHotLiveRoomJoinStat(week), nil
	})
}

func ListRecentWeeklyHotLiveRoomJoinStats(limit int) []*entity.WeeklyHotLiveRoomJoinStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyHotLiveRoomJoinStatsFromDB(limit)
}

func loadRecentWeeklyHotLiveRoomJoinStatsFromDB(limit int) []*entity.WeeklyHotLiveRoomJoinStat {
	list := make([]*entity.WeeklyHotLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyHotLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyHotLiveRoomJoinStats(list)
	return list
}

func reverseWeeklyHotLiveRoomJoinStats(list []*entity.WeeklyHotLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyHotLiveRoomJoinStat(data *entity.WeeklyHotLiveRoomJoinStat) {
	if data == nil || data.ID == "" || weeklyHotLiveRoomJoinStatCacheMgr == nil {
		return
	}
	weeklyHotLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
