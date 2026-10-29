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
	weeklyHotLiveRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.WeeklyHotLiveRoomGiftClickStat]
)

func initWeeklyHotLiveRoomGiftClickStatDao() {
	weeklyHotLiveRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyHotLiveRoomGiftClickStat]()
}

func GetWeeklyHotLiveRoomGiftClickStatByWeek(week string) *entity.WeeklyHotLiveRoomGiftClickStat {
	return weeklyHotLiveRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyHotLiveRoomGiftClickStat, error) {
		var data *entity.WeeklyHotLiveRoomGiftClickStat
		_ = g.Model(string(entity.TbWeeklyHotLiveRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyHotLiveRoomGiftClickStat(week), nil
	})
}

func ListRecentWeeklyHotLiveRoomGiftClickStats(limit int) []*entity.WeeklyHotLiveRoomGiftClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyHotLiveRoomGiftClickStatsFromDB(limit)
}

func loadRecentWeeklyHotLiveRoomGiftClickStatsFromDB(limit int) []*entity.WeeklyHotLiveRoomGiftClickStat {
	list := make([]*entity.WeeklyHotLiveRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyHotLiveRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyHotLiveRoomGiftClickStats(list)
	return list
}

func reverseWeeklyHotLiveRoomGiftClickStats(list []*entity.WeeklyHotLiveRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyHotLiveRoomGiftClickStat(data *entity.WeeklyHotLiveRoomGiftClickStat) {
	if data == nil || data.ID == "" || weeklyHotLiveRoomGiftClickStatCacheMgr == nil {
		return
	}
	weeklyHotLiveRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
