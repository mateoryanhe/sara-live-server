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
	dailyMiniGameRoundStartStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameRoundStartStat]
	dailyMiniGameRoundStartStatListCacheMgr *cache.ListCache[*entity.DailyMiniGameRoundStartStat]
)

func initDailyMiniGameRoundStartStatDao() {
	dailyMiniGameRoundStartStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameRoundStartStat]()
	dailyMiniGameRoundStartStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyMiniGameRoundStartStat]()
}

func GetDailyMiniGameRoundStartStatByDate(date string) *entity.DailyMiniGameRoundStartStat {
	return dailyMiniGameRoundStartStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameRoundStartStat, error) {
		var data *entity.DailyMiniGameRoundStartStat
		_ = g.Model(string(entity.TbDailyMiniGameRoundStartStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameRoundStartStat(date), nil
	})
}

func ListRecentDailyMiniGameRoundStartStats(limit int) []*entity.DailyMiniGameRoundStartStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyMiniGameRoundStartStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyMiniGameRoundStartStat, error) {
		return loadRecentDailyMiniGameRoundStartStatsFromDB(limit), nil
	})
}

func loadRecentDailyMiniGameRoundStartStatsFromDB(limit int) []*entity.DailyMiniGameRoundStartStat {
	list := make([]*entity.DailyMiniGameRoundStartStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameRoundStartStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameRoundStartStats(list)
	return list
}

func reverseDailyMiniGameRoundStartStats(list []*entity.DailyMiniGameRoundStartStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameRoundStartStat(data *entity.DailyMiniGameRoundStartStat) {
	if data == nil || data.ID == "" || dailyMiniGameRoundStartStatCacheMgr == nil {
		return
	}
	dailyMiniGameRoundStartStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
