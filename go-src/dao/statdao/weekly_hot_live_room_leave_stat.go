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
	weeklyHotLiveRoomLeaveStatCacheMgr     *cache.RowCache[*entity.WeeklyHotLiveRoomLeaveStat]
	weeklyHotLiveRoomLeaveStatListCacheMgr *cache.ListCache[*entity.WeeklyHotLiveRoomLeaveStat]
)

func initWeeklyHotLiveRoomLeaveStatDao() {
	weeklyHotLiveRoomLeaveStatCacheMgr = cache.NewRowCache[*entity.WeeklyHotLiveRoomLeaveStat]()
	weeklyHotLiveRoomLeaveStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyHotLiveRoomLeaveStat]()
}

func GetWeeklyHotLiveRoomLeaveStatByWeek(week string) *entity.WeeklyHotLiveRoomLeaveStat {
	return weeklyHotLiveRoomLeaveStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyHotLiveRoomLeaveStat, error) {
		var data *entity.WeeklyHotLiveRoomLeaveStat
		_ = g.Model(string(entity.TbWeeklyHotLiveRoomLeaveStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyHotLiveRoomLeaveStat(week), nil
	})
}

func ListRecentWeeklyHotLiveRoomLeaveStats(limit int) []*entity.WeeklyHotLiveRoomLeaveStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyHotLiveRoomLeaveStatKey(time.Now()))
	return weeklyHotLiveRoomLeaveStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyHotLiveRoomLeaveStat, error) {
		return loadRecentWeeklyHotLiveRoomLeaveStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyHotLiveRoomLeaveStatsFromDB(limit int) []*entity.WeeklyHotLiveRoomLeaveStat {
	list := make([]*entity.WeeklyHotLiveRoomLeaveStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyHotLiveRoomLeaveStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyHotLiveRoomLeaveStats(list)
	return list
}

func reverseWeeklyHotLiveRoomLeaveStats(list []*entity.WeeklyHotLiveRoomLeaveStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyHotLiveRoomLeaveStat(data *entity.WeeklyHotLiveRoomLeaveStat) {
	if data == nil || data.ID == "" || weeklyHotLiveRoomLeaveStatCacheMgr == nil {
		return
	}
	weeklyHotLiveRoomLeaveStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
