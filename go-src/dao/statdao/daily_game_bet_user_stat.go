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
	dailyGameBetUserStatCacheMgr     *cache.RowCache[*entity.DailyGameBetUserStat]
	dailyGameBetUserStatListCacheMgr *cache.ListCache[*entity.DailyGameBetUserStat]
)

func initDailyGameBetUserStatDao() {
	dailyGameBetUserStatCacheMgr = cache.NewRowCache[*entity.DailyGameBetUserStat]()
	dailyGameBetUserStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyGameBetUserStat]()
}

func GetDailyGameBetUserStatByDate(date string) *entity.DailyGameBetUserStat {
	return dailyGameBetUserStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyGameBetUserStat, error) {
		var data *entity.DailyGameBetUserStat
		_ = g.Model(string(entity.TbDailyGameBetUserStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyGameBetUserStat(date), nil
	})
}

func ListRecentDailyGameBetUserStats(limit int) []*entity.DailyGameBetUserStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyGameBetUserStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyGameBetUserStat, error) {
		return loadRecentDailyGameBetUserStatsFromDB(limit), nil
	})
}

func loadRecentDailyGameBetUserStatsFromDB(limit int) []*entity.DailyGameBetUserStat {
	list := make([]*entity.DailyGameBetUserStat, 0, limit)
	_ = g.Model(string(entity.TbDailyGameBetUserStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyGameBetUserStats(list)
	return list
}

func reverseDailyGameBetUserStats(list []*entity.DailyGameBetUserStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyGameBetUserStat(data *entity.DailyGameBetUserStat) {
	if data == nil || data.ID == "" || dailyGameBetUserStatCacheMgr == nil {
		return
	}
	dailyGameBetUserStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
