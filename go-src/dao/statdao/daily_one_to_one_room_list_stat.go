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
	dailyOneToOneRoomListStatCacheMgr     *cache.RowCache[*entity.DailyOneToOneRoomListStat]
	dailyOneToOneRoomListStatListCacheMgr *cache.ListCache[*entity.DailyOneToOneRoomListStat]
)

func initDailyOneToOneRoomListStatDao() {
	dailyOneToOneRoomListStatCacheMgr = cache.NewRowCache[*entity.DailyOneToOneRoomListStat]()
	dailyOneToOneRoomListStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyOneToOneRoomListStat]()
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
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyOneToOneRoomListStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyOneToOneRoomListStat, error) {
		return loadRecentDailyOneToOneRoomListStatsFromDB(limit), nil
	})
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
