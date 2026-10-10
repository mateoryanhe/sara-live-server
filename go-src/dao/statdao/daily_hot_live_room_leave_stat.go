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
	dailyHotLiveRoomLeaveStatCacheMgr     *cache.RowCache[*entity.DailyHotLiveRoomLeaveStat]
	dailyHotLiveRoomLeaveStatListCacheMgr *cache.ListCache[*entity.DailyHotLiveRoomLeaveStat]
)

func initDailyHotLiveRoomLeaveStatDao() {
	dailyHotLiveRoomLeaveStatCacheMgr = cache.NewRowCache[*entity.DailyHotLiveRoomLeaveStat]()
	dailyHotLiveRoomLeaveStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyHotLiveRoomLeaveStat]()
}

func GetDailyHotLiveRoomLeaveStatByDate(date string) *entity.DailyHotLiveRoomLeaveStat {
	return dailyHotLiveRoomLeaveStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyHotLiveRoomLeaveStat, error) {
		var data *entity.DailyHotLiveRoomLeaveStat
		_ = g.Model(string(entity.TbDailyHotLiveRoomLeaveStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyHotLiveRoomLeaveStat(date), nil
	})
}

func ListRecentDailyHotLiveRoomLeaveStats(limit int) []*entity.DailyHotLiveRoomLeaveStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyHotLiveRoomLeaveStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyHotLiveRoomLeaveStat, error) {
		return loadRecentDailyHotLiveRoomLeaveStatsFromDB(limit), nil
	})
}

func loadRecentDailyHotLiveRoomLeaveStatsFromDB(limit int) []*entity.DailyHotLiveRoomLeaveStat {
	list := make([]*entity.DailyHotLiveRoomLeaveStat, 0, limit)
	_ = g.Model(string(entity.TbDailyHotLiveRoomLeaveStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyHotLiveRoomLeaveStats(list)
	return list
}

func reverseDailyHotLiveRoomLeaveStats(list []*entity.DailyHotLiveRoomLeaveStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyHotLiveRoomLeaveStat(data *entity.DailyHotLiveRoomLeaveStat) {
	if data == nil || data.ID == "" || dailyHotLiveRoomLeaveStatCacheMgr == nil {
		return
	}
	dailyHotLiveRoomLeaveStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
