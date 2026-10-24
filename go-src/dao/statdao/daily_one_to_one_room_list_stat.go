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
	dailyOneToOneRoomListStatCacheMgr *cache.RowCache[*entity.DailyOneToOneRoomListStat]
)

func initDailyOneToOneRoomListStatDao() {
	dailyOneToOneRoomListStatCacheMgr = cache.NewRowCache[*entity.DailyOneToOneRoomListStat]()
}

func GetDailyOneToOneRoomListStatByDate(date string) *entity.DailyOneToOneRoomListStat {
	return dailyOneToOneRoomListStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyOneToOneRoomListStat, error) {
		var data *entity.DailyOneToOneRoomListStat
		_ = g.Model(string(entity.TbDailyOneToOneRoomListStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyOneToOneRoomListStat(date), nil
	})
}

func ListRecentDailyOneToOneRoomListStats(limit int) []*entity.DailyOneToOneRoomListStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyOneToOneRoomListStatsFromDB(limit)
}

func loadRecentDailyOneToOneRoomListStatsFromDB(limit int) []*entity.DailyOneToOneRoomListStat {
	list := make([]*entity.DailyOneToOneRoomListStat, 0, limit)
	_ = g.Model(string(entity.TbDailyOneToOneRoomListStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyOneToOneRoomListStats(list)
	return list
}

func reverseDailyOneToOneRoomListStats(list []*entity.DailyOneToOneRoomListStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyOneToOneRoomListStat(data *entity.DailyOneToOneRoomListStat) {
	if data == nil || data.ID == "" || dailyOneToOneRoomListStatCacheMgr == nil {
		return
	}
	dailyOneToOneRoomListStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
