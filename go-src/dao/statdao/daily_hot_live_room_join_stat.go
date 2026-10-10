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
	dailyHotLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.DailyHotLiveRoomJoinStat]
	dailyHotLiveRoomJoinStatListCacheMgr *cache.ListCache[*entity.DailyHotLiveRoomJoinStat]
)

func initDailyHotLiveRoomJoinStatDao() {
	dailyHotLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.DailyHotLiveRoomJoinStat]()
	dailyHotLiveRoomJoinStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyHotLiveRoomJoinStat]()
}

func GetDailyHotLiveRoomJoinStatByDate(date string) *entity.DailyHotLiveRoomJoinStat {
	return dailyHotLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyHotLiveRoomJoinStat, error) {
		var data *entity.DailyHotLiveRoomJoinStat
		_ = g.Model(string(entity.TbDailyHotLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyHotLiveRoomJoinStat(date), nil
	})
}

func ListRecentDailyHotLiveRoomJoinStats(limit int) []*entity.DailyHotLiveRoomJoinStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyHotLiveRoomJoinStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyHotLiveRoomJoinStat, error) {
		return loadRecentDailyHotLiveRoomJoinStatsFromDB(limit), nil
	})
}

func loadRecentDailyHotLiveRoomJoinStatsFromDB(limit int) []*entity.DailyHotLiveRoomJoinStat {
	list := make([]*entity.DailyHotLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbDailyHotLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyHotLiveRoomJoinStats(list)
	return list
}

func reverseDailyHotLiveRoomJoinStats(list []*entity.DailyHotLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyHotLiveRoomJoinStat(data *entity.DailyHotLiveRoomJoinStat) {
	if data == nil || data.ID == "" || dailyHotLiveRoomJoinStatCacheMgr == nil {
		return
	}
	dailyHotLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
