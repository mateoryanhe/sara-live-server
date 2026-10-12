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
	dailyChargePanelEnterStatCacheMgr     *cache.RowCache[*entity.DailyChargePanelEnterStat]
	dailyChargePanelEnterStatListCacheMgr *cache.ListCache[*entity.DailyChargePanelEnterStat]
)

func initDailyChargePanelEnterStatDao() {
	dailyChargePanelEnterStatCacheMgr = cache.NewRowCache[*entity.DailyChargePanelEnterStat]()
	dailyChargePanelEnterStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyChargePanelEnterStat]()
}

func GetDailyChargePanelEnterStatByDate(date string) *entity.DailyChargePanelEnterStat {
	return dailyChargePanelEnterStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyChargePanelEnterStat, error) {
		var data *entity.DailyChargePanelEnterStat
		_ = g.Model(string(entity.TbDailyChargePanelEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyChargePanelEnterStat(date), nil
	})
}

func ListRecentDailyChargePanelEnterStats(limit int) []*entity.DailyChargePanelEnterStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyChargePanelEnterStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyChargePanelEnterStat, error) {
		return loadRecentDailyChargePanelEnterStatsFromDB(limit), nil
	})
}

func loadRecentDailyChargePanelEnterStatsFromDB(limit int) []*entity.DailyChargePanelEnterStat {
	list := make([]*entity.DailyChargePanelEnterStat, 0, limit)
	_ = g.Model(string(entity.TbDailyChargePanelEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyChargePanelEnterStats(list)
	return list
}

func reverseDailyChargePanelEnterStats(list []*entity.DailyChargePanelEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyChargePanelEnterStat(data *entity.DailyChargePanelEnterStat) {
	if data == nil || data.ID == "" || dailyChargePanelEnterStatCacheMgr == nil {
		return
	}
	dailyChargePanelEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
