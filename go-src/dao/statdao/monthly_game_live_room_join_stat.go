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
	monthlyGameLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.MonthlyGameLiveRoomJoinStat]
	monthlyGameLiveRoomJoinStatListCacheMgr *cache.ListCache[*entity.MonthlyGameLiveRoomJoinStat]
)

func initMonthlyGameLiveRoomJoinStatDao() {
	monthlyGameLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.MonthlyGameLiveRoomJoinStat]()
	monthlyGameLiveRoomJoinStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyGameLiveRoomJoinStat]()
}

func GetMonthlyGameLiveRoomJoinStatByMonth(month string) *entity.MonthlyGameLiveRoomJoinStat {
	return monthlyGameLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyGameLiveRoomJoinStat, error) {
		var data *entity.MonthlyGameLiveRoomJoinStat
		_ = g.Model(string(entity.TbMonthlyGameLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyGameLiveRoomJoinStat(month), nil
	})
}

func ListRecentMonthlyGameLiveRoomJoinStats(limit int) []*entity.MonthlyGameLiveRoomJoinStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyGameLiveRoomJoinStatKey(time.Now()))
	return monthlyGameLiveRoomJoinStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyGameLiveRoomJoinStat, error) {
		return loadRecentMonthlyGameLiveRoomJoinStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyGameLiveRoomJoinStatsFromDB(limit int) []*entity.MonthlyGameLiveRoomJoinStat {
	list := make([]*entity.MonthlyGameLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyGameLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyGameLiveRoomJoinStats(list)
	return list
}

func reverseMonthlyGameLiveRoomJoinStats(list []*entity.MonthlyGameLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyGameLiveRoomJoinStat(data *entity.MonthlyGameLiveRoomJoinStat) {
	if data == nil || data.ID == "" || monthlyGameLiveRoomJoinStatCacheMgr == nil {
		return
	}
	monthlyGameLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
