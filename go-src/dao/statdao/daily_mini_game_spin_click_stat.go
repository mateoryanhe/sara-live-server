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
	dailyMiniGameSpinClickStatCacheMgr     *cache.RowCache[*entity.DailyMiniGameSpinClickStat]
	dailyMiniGameSpinClickStatListCacheMgr *cache.ListCache[*entity.DailyMiniGameSpinClickStat]
)

func initDailyMiniGameSpinClickStatDao() {
	dailyMiniGameSpinClickStatCacheMgr = cache.NewRowCache[*entity.DailyMiniGameSpinClickStat]()
	dailyMiniGameSpinClickStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyMiniGameSpinClickStat]()
}

func GetDailyMiniGameSpinClickStatByDate(date string) *entity.DailyMiniGameSpinClickStat {
	return dailyMiniGameSpinClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyMiniGameSpinClickStat, error) {
		var data *entity.DailyMiniGameSpinClickStat
		_ = g.Model(string(entity.TbDailyMiniGameSpinClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyMiniGameSpinClickStat(date), nil
	})
}

func ListRecentDailyMiniGameSpinClickStats(limit int) []*entity.DailyMiniGameSpinClickStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyMiniGameSpinClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyMiniGameSpinClickStat, error) {
		return loadRecentDailyMiniGameSpinClickStatsFromDB(limit), nil
	})
}

func loadRecentDailyMiniGameSpinClickStatsFromDB(limit int) []*entity.DailyMiniGameSpinClickStat {
	list := make([]*entity.DailyMiniGameSpinClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyMiniGameSpinClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyMiniGameSpinClickStats(list)
	return list
}

func reverseDailyMiniGameSpinClickStats(list []*entity.DailyMiniGameSpinClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyMiniGameSpinClickStat(data *entity.DailyMiniGameSpinClickStat) {
	if data == nil || data.ID == "" || dailyMiniGameSpinClickStatCacheMgr == nil {
		return
	}
	dailyMiniGameSpinClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
