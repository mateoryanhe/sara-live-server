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
	monthlyOneToOneRoomListStatCacheMgr     *cache.RowCache[*entity.MonthlyOneToOneRoomListStat]
)

func initMonthlyOneToOneRoomListStatDao() {
	monthlyOneToOneRoomListStatCacheMgr = cache.NewRowCache[*entity.MonthlyOneToOneRoomListStat]()
}

func GetMonthlyOneToOneRoomListStatByMonth(month string) *entity.MonthlyOneToOneRoomListStat {
	return monthlyOneToOneRoomListStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyOneToOneRoomListStat, error) {
		var data *entity.MonthlyOneToOneRoomListStat
		_ = g.Model(string(entity.TbMonthlyOneToOneRoomListStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyOneToOneRoomListStat(month), nil
	})
}

func ListRecentMonthlyOneToOneRoomListStats(limit int) []*entity.MonthlyOneToOneRoomListStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentMonthlyOneToOneRoomListStatsFromDB(limit)
}

func loadRecentMonthlyOneToOneRoomListStatsFromDB(limit int) []*entity.MonthlyOneToOneRoomListStat {
	list := make([]*entity.MonthlyOneToOneRoomListStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyOneToOneRoomListStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyOneToOneRoomListStats(list)
	return list
}

func reverseMonthlyOneToOneRoomListStats(list []*entity.MonthlyOneToOneRoomListStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyOneToOneRoomListStat(data *entity.MonthlyOneToOneRoomListStat) {
	if data == nil || data.ID == "" || monthlyOneToOneRoomListStatCacheMgr == nil {
		return
	}
	monthlyOneToOneRoomListStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
