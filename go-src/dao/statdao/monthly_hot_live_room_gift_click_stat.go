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
	monthlyHotLiveRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.MonthlyHotLiveRoomGiftClickStat]
	monthlyHotLiveRoomGiftClickStatListCacheMgr *cache.ListCache[*entity.MonthlyHotLiveRoomGiftClickStat]
)

func initMonthlyHotLiveRoomGiftClickStatDao() {
	monthlyHotLiveRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyHotLiveRoomGiftClickStat]()
	monthlyHotLiveRoomGiftClickStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyHotLiveRoomGiftClickStat]()
}

func GetMonthlyHotLiveRoomGiftClickStatByMonth(month string) *entity.MonthlyHotLiveRoomGiftClickStat {
	return monthlyHotLiveRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyHotLiveRoomGiftClickStat, error) {
		var data *entity.MonthlyHotLiveRoomGiftClickStat
		_ = g.Model(string(entity.TbMonthlyHotLiveRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyHotLiveRoomGiftClickStat(month), nil
	})
}

func ListRecentMonthlyHotLiveRoomGiftClickStats(limit int) []*entity.MonthlyHotLiveRoomGiftClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyHotLiveRoomGiftClickStatKey(time.Now()))
	return monthlyHotLiveRoomGiftClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyHotLiveRoomGiftClickStat, error) {
		return loadRecentMonthlyHotLiveRoomGiftClickStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyHotLiveRoomGiftClickStatsFromDB(limit int) []*entity.MonthlyHotLiveRoomGiftClickStat {
	list := make([]*entity.MonthlyHotLiveRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyHotLiveRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyHotLiveRoomGiftClickStats(list)
	return list
}

func reverseMonthlyHotLiveRoomGiftClickStats(list []*entity.MonthlyHotLiveRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyHotLiveRoomGiftClickStat(data *entity.MonthlyHotLiveRoomGiftClickStat) {
	if data == nil || data.ID == "" || monthlyHotLiveRoomGiftClickStatCacheMgr == nil {
		return
	}
	monthlyHotLiveRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
