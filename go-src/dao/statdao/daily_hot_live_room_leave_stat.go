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
	dailyHotLiveRoomLeaveStatCacheMgr     *cache.RowCache[*entity.DailyHotLiveRoomLeaveStat]
)

func initDailyHotLiveRoomLeaveStatDao() {
	dailyHotLiveRoomLeaveStatCacheMgr = cache.NewRowCache[*entity.DailyHotLiveRoomLeaveStat]()
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
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyHotLiveRoomLeaveStatsFromDB(limit)
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
