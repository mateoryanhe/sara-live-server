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
	dailyChargeResultServerStatCacheMgr     *cache.RowCache[*entity.DailyChargeResultServerStat]
	dailyChargeResultServerStatListCacheMgr *cache.ListCache[*entity.DailyChargeResultServerStat]
)

func initDailyChargeResultServerStatDao() {
	dailyChargeResultServerStatCacheMgr = cache.NewRowCache[*entity.DailyChargeResultServerStat]()
	dailyChargeResultServerStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyChargeResultServerStat]()
}

func GetDailyChargeResultServerStatByDate(date string) *entity.DailyChargeResultServerStat {
	return dailyChargeResultServerStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyChargeResultServerStat, error) {
		var data *entity.DailyChargeResultServerStat
		_ = g.Model(string(entity.TbDailyChargeResultServerStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyChargeResultServerStat(date), nil
	})
}

func ListRecentDailyChargeResultServerStats(limit int) []*entity.DailyChargeResultServerStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyChargeResultServerStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyChargeResultServerStat, error) {
		return loadRecentDailyChargeResultServerStatsFromDB(limit), nil
	})
}

func loadRecentDailyChargeResultServerStatsFromDB(limit int) []*entity.DailyChargeResultServerStat {
	list := make([]*entity.DailyChargeResultServerStat, 0, limit)
	_ = g.Model(string(entity.TbDailyChargeResultServerStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyChargeResultServerStats(list)
	return list
}

func reverseDailyChargeResultServerStats(list []*entity.DailyChargeResultServerStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyChargeResultServerStat(data *entity.DailyChargeResultServerStat) {
	if data == nil || data.ID == "" || dailyChargeResultServerStatCacheMgr == nil {
		return
	}
	dailyChargeResultServerStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
