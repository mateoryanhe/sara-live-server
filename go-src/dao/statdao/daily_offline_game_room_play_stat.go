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
	dailyOfflineGameRoomPlayStatCacheMgr     *cache.RowCache[*entity.DailyOfflineGameRoomPlayStat]
	dailyOfflineGameRoomPlayStatListCacheMgr *cache.ListCache[*entity.DailyOfflineGameRoomPlayStat]
)

func initDailyOfflineGameRoomPlayStatDao() {
	dailyOfflineGameRoomPlayStatCacheMgr = cache.NewRowCache[*entity.DailyOfflineGameRoomPlayStat]()
	dailyOfflineGameRoomPlayStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyOfflineGameRoomPlayStat]()
}

func GetDailyOfflineGameRoomPlayStatByDate(date string) *entity.DailyOfflineGameRoomPlayStat {
	return dailyOfflineGameRoomPlayStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyOfflineGameRoomPlayStat, error) {
		var data *entity.DailyOfflineGameRoomPlayStat
		_ = g.Model(string(entity.TbDailyOfflineGameRoomPlayStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyOfflineGameRoomPlayStat(date), nil
	})
}

func ListRecentDailyOfflineGameRoomPlayStats(limit int) []*entity.DailyOfflineGameRoomPlayStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyOfflineGameRoomPlayStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyOfflineGameRoomPlayStat, error) {
		return loadRecentDailyOfflineGameRoomPlayStatsFromDB(limit), nil
	})
}

func loadRecentDailyOfflineGameRoomPlayStatsFromDB(limit int) []*entity.DailyOfflineGameRoomPlayStat {
	list := make([]*entity.DailyOfflineGameRoomPlayStat, 0, limit)
	_ = g.Model(string(entity.TbDailyOfflineGameRoomPlayStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyOfflineGameRoomPlayStats(list)
	return list
}

func reverseDailyOfflineGameRoomPlayStats(list []*entity.DailyOfflineGameRoomPlayStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyOfflineGameRoomPlayStat(data *entity.DailyOfflineGameRoomPlayStat) {
	if data == nil || data.ID == "" || dailyOfflineGameRoomPlayStatCacheMgr == nil {
		return
	}
	dailyOfflineGameRoomPlayStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
