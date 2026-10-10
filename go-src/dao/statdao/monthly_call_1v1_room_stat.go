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
	monthlyCall1v1RoomStatCacheMgr     *cache.RowCache[*entity.MonthlyCall1v1RoomStat]
	monthlyCall1v1RoomStatListCacheMgr *cache.ListCache[*entity.MonthlyCall1v1RoomStat]
)

func initMonthlyCall1v1RoomStatDao() {
	monthlyCall1v1RoomStatCacheMgr = cache.NewRowCache[*entity.MonthlyCall1v1RoomStat]()
	monthlyCall1v1RoomStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyCall1v1RoomStat]()
}

func GetMonthlyCall1v1RoomStatByMonth(month string) *entity.MonthlyCall1v1RoomStat {
	return monthlyCall1v1RoomStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyCall1v1RoomStat, error) {
		var data *entity.MonthlyCall1v1RoomStat
		_ = g.Model(string(entity.TbMonthlyCall1v1RoomStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyCall1v1RoomStat(month), nil
	})
}

func ListRecentMonthlyCall1v1RoomStats(limit int) []*entity.MonthlyCall1v1RoomStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyCall1v1RoomStatKey(time.Now()))
	return monthlyCall1v1RoomStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyCall1v1RoomStat, error) {
		return loadRecentMonthlyCall1v1RoomStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyCall1v1RoomStatsFromDB(limit int) []*entity.MonthlyCall1v1RoomStat {
	list := make([]*entity.MonthlyCall1v1RoomStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyCall1v1RoomStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyCall1v1RoomStats(list)
	return list
}

func reverseMonthlyCall1v1RoomStats(list []*entity.MonthlyCall1v1RoomStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyCall1v1RoomStat(data *entity.MonthlyCall1v1RoomStat) {
	if data == nil || data.ID == "" || monthlyCall1v1RoomStatCacheMgr == nil {
		return
	}
	monthlyCall1v1RoomStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
