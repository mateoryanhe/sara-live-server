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
	dailyCall1v1RoomStatCacheMgr     *cache.RowCache[*entity.DailyCall1v1RoomStat]
	dailyCall1v1RoomStatListCacheMgr *cache.ListCache[*entity.DailyCall1v1RoomStat]
)

func initDailyCall1v1RoomStatDao() {
	dailyCall1v1RoomStatCacheMgr = cache.NewRowCache[*entity.DailyCall1v1RoomStat]()
	dailyCall1v1RoomStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyCall1v1RoomStat]()
}

func GetDailyCall1v1RoomStatByDate(date string) *entity.DailyCall1v1RoomStat {
	return dailyCall1v1RoomStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyCall1v1RoomStat, error) {
		var data *entity.DailyCall1v1RoomStat
		_ = g.Model(string(entity.TbDailyCall1v1RoomStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyCall1v1RoomStat(date), nil
	})
}

func ListRecentDailyCall1v1RoomStats(limit int) []*entity.DailyCall1v1RoomStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyCall1v1RoomStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyCall1v1RoomStat, error) {
		return loadRecentDailyCall1v1RoomStatsFromDB(limit), nil
	})
}

func loadRecentDailyCall1v1RoomStatsFromDB(limit int) []*entity.DailyCall1v1RoomStat {
	list := make([]*entity.DailyCall1v1RoomStat, 0, limit)
	_ = g.Model(string(entity.TbDailyCall1v1RoomStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyCall1v1RoomStats(list)
	return list
}

func reverseDailyCall1v1RoomStats(list []*entity.DailyCall1v1RoomStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyCall1v1RoomStat(data *entity.DailyCall1v1RoomStat) {
	if data == nil || data.ID == "" || dailyCall1v1RoomStatCacheMgr == nil {
		return
	}
	dailyCall1v1RoomStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
