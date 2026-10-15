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
	monthlyOfflineGameRoomPlayStatCacheMgr     *cache.RowCache[*entity.MonthlyOfflineGameRoomPlayStat]
	monthlyOfflineGameRoomPlayStatListCacheMgr *cache.ListCache[*entity.MonthlyOfflineGameRoomPlayStat]
)

func initMonthlyOfflineGameRoomPlayStatDao() {
	monthlyOfflineGameRoomPlayStatCacheMgr = cache.NewRowCache[*entity.MonthlyOfflineGameRoomPlayStat]()
	monthlyOfflineGameRoomPlayStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyOfflineGameRoomPlayStat]()
}

func GetMonthlyOfflineGameRoomPlayStatByMonth(month string) *entity.MonthlyOfflineGameRoomPlayStat {
	return monthlyOfflineGameRoomPlayStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyOfflineGameRoomPlayStat, error) {
		var data *entity.MonthlyOfflineGameRoomPlayStat
		_ = g.Model(string(entity.TbMonthlyOfflineGameRoomPlayStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyOfflineGameRoomPlayStat(month), nil
	})
}

func ListRecentMonthlyOfflineGameRoomPlayStats(limit int) []*entity.MonthlyOfflineGameRoomPlayStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyOfflineGameRoomPlayStatKey(time.Now()))
	return monthlyOfflineGameRoomPlayStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyOfflineGameRoomPlayStat, error) {
		return loadRecentMonthlyOfflineGameRoomPlayStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyOfflineGameRoomPlayStatsFromDB(limit int) []*entity.MonthlyOfflineGameRoomPlayStat {
	list := make([]*entity.MonthlyOfflineGameRoomPlayStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyOfflineGameRoomPlayStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyOfflineGameRoomPlayStats(list)
	return list
}

func reverseMonthlyOfflineGameRoomPlayStats(list []*entity.MonthlyOfflineGameRoomPlayStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyOfflineGameRoomPlayStat(data *entity.MonthlyOfflineGameRoomPlayStat) {
	if data == nil || data.ID == "" || monthlyOfflineGameRoomPlayStatCacheMgr == nil {
		return
	}
	monthlyOfflineGameRoomPlayStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
