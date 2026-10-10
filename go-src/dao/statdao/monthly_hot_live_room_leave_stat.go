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
	monthlyHotLiveRoomLeaveStatCacheMgr     *cache.RowCache[*entity.MonthlyHotLiveRoomLeaveStat]
	monthlyHotLiveRoomLeaveStatListCacheMgr *cache.ListCache[*entity.MonthlyHotLiveRoomLeaveStat]
)

func initMonthlyHotLiveRoomLeaveStatDao() {
	monthlyHotLiveRoomLeaveStatCacheMgr = cache.NewRowCache[*entity.MonthlyHotLiveRoomLeaveStat]()
	monthlyHotLiveRoomLeaveStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyHotLiveRoomLeaveStat]()
}

func GetMonthlyHotLiveRoomLeaveStatByMonth(month string) *entity.MonthlyHotLiveRoomLeaveStat {
	return monthlyHotLiveRoomLeaveStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyHotLiveRoomLeaveStat, error) {
		var data *entity.MonthlyHotLiveRoomLeaveStat
		_ = g.Model(string(entity.TbMonthlyHotLiveRoomLeaveStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyHotLiveRoomLeaveStat(month), nil
	})
}

func ListRecentMonthlyHotLiveRoomLeaveStats(limit int) []*entity.MonthlyHotLiveRoomLeaveStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyHotLiveRoomLeaveStatKey(time.Now()))
	return monthlyHotLiveRoomLeaveStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyHotLiveRoomLeaveStat, error) {
		return loadRecentMonthlyHotLiveRoomLeaveStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyHotLiveRoomLeaveStatsFromDB(limit int) []*entity.MonthlyHotLiveRoomLeaveStat {
	list := make([]*entity.MonthlyHotLiveRoomLeaveStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyHotLiveRoomLeaveStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyHotLiveRoomLeaveStats(list)
	return list
}

func reverseMonthlyHotLiveRoomLeaveStats(list []*entity.MonthlyHotLiveRoomLeaveStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyHotLiveRoomLeaveStat(data *entity.MonthlyHotLiveRoomLeaveStat) {
	if data == nil || data.ID == "" || monthlyHotLiveRoomLeaveStatCacheMgr == nil {
		return
	}
	monthlyHotLiveRoomLeaveStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
