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
	dailyGameLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.DailyGameLiveRoomJoinStat]
)

func initDailyGameLiveRoomJoinStatDao() {
	dailyGameLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.DailyGameLiveRoomJoinStat]()
}

func GetDailyGameLiveRoomJoinStatByDate(date string) *entity.DailyGameLiveRoomJoinStat {
	return dailyGameLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGameLiveRoomJoinStat, error) {
		var data *entity.DailyGameLiveRoomJoinStat
		_ = g.Model(string(entity.TbDailyGameLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGameLiveRoomJoinStat(date), nil
	})
}

func ListRecentDailyGameLiveRoomJoinStats(limit int) []*entity.DailyGameLiveRoomJoinStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyGameLiveRoomJoinStatsFromDB(limit)
}

func loadRecentDailyGameLiveRoomJoinStatsFromDB(limit int) []*entity.DailyGameLiveRoomJoinStat {
	list := make([]*entity.DailyGameLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGameLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGameLiveRoomJoinStats(list)
	return list
}

func reverseDailyGameLiveRoomJoinStats(list []*entity.DailyGameLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGameLiveRoomJoinStat(data *entity.DailyGameLiveRoomJoinStat) {
	if data == nil || data.ID == "" || dailyGameLiveRoomJoinStatCacheMgr == nil {
		return
	}
	dailyGameLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
