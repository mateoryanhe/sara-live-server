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
	monthlyHotLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.MonthlyHotLiveRoomJoinStat]
	monthlyHotLiveRoomJoinStatListCacheMgr *cache.ListCache[*entity.MonthlyHotLiveRoomJoinStat]
)

func initMonthlyHotLiveRoomJoinStatDao() {
	monthlyHotLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.MonthlyHotLiveRoomJoinStat]()
	monthlyHotLiveRoomJoinStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyHotLiveRoomJoinStat]()
}

func GetMonthlyHotLiveRoomJoinStatByMonth(month string) *entity.MonthlyHotLiveRoomJoinStat {
	return monthlyHotLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyHotLiveRoomJoinStat, error) {
		var data *entity.MonthlyHotLiveRoomJoinStat
		_ = g.Model(string(entity.TbMonthlyHotLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyHotLiveRoomJoinStat(month), nil
	})
}

func ListRecentMonthlyHotLiveRoomJoinStats(limit int) []*entity.MonthlyHotLiveRoomJoinStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyHotLiveRoomJoinStatKey(time.Now()))
	return monthlyHotLiveRoomJoinStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyHotLiveRoomJoinStat, error) {
		return loadRecentMonthlyHotLiveRoomJoinStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyHotLiveRoomJoinStatsFromDB(limit int) []*entity.MonthlyHotLiveRoomJoinStat {
	list := make([]*entity.MonthlyHotLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyHotLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyHotLiveRoomJoinStats(list)
	return list
}

func reverseMonthlyHotLiveRoomJoinStats(list []*entity.MonthlyHotLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyHotLiveRoomJoinStat(data *entity.MonthlyHotLiveRoomJoinStat) {
	if data == nil || data.ID == "" || monthlyHotLiveRoomJoinStatCacheMgr == nil {
		return
	}
	monthlyHotLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
