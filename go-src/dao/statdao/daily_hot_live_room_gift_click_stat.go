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
	dailyHotLiveRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.DailyHotLiveRoomGiftClickStat]
)

func initDailyHotLiveRoomGiftClickStatDao() {
	dailyHotLiveRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.DailyHotLiveRoomGiftClickStat]()
}

func GetDailyHotLiveRoomGiftClickStatByDate(date string) *entity.DailyHotLiveRoomGiftClickStat {
	return dailyHotLiveRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyHotLiveRoomGiftClickStat, error) {
		var data *entity.DailyHotLiveRoomGiftClickStat
		_ = g.Model(string(entity.TbDailyHotLiveRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyHotLiveRoomGiftClickStat(date), nil
	})
}

func ListRecentDailyHotLiveRoomGiftClickStats(limit int) []*entity.DailyHotLiveRoomGiftClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyHotLiveRoomGiftClickStatsFromDB(limit)
}

func loadRecentDailyHotLiveRoomGiftClickStatsFromDB(limit int) []*entity.DailyHotLiveRoomGiftClickStat {
	list := make([]*entity.DailyHotLiveRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyHotLiveRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyHotLiveRoomGiftClickStats(list)
	return list
}

func reverseDailyHotLiveRoomGiftClickStats(list []*entity.DailyHotLiveRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyHotLiveRoomGiftClickStat(data *entity.DailyHotLiveRoomGiftClickStat) {
	if data == nil || data.ID == "" || dailyHotLiveRoomGiftClickStatCacheMgr == nil {
		return
	}
	dailyHotLiveRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
