package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyEmailAccountStatCacheMgr *cache.RowCache[*entity.DailyEmailAccountStat]

func initDailyEmailAccountStatDao() {
	dailyEmailAccountStatCacheMgr = cache.NewRowCache[*entity.DailyEmailAccountStat]()
}

func GetDailyEmailAccountStatByDate(date string) *entity.DailyEmailAccountStat {
	return dailyEmailAccountStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyEmailAccountStat, error) {
		var data *entity.DailyEmailAccountStat
		_ = g.Model(string(entity.TbDailyEmailAccountStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyEmailAccountStat(date), nil
	})
}

func ListRecentDailyEmailAccountStats(limit int) []*entity.DailyEmailAccountStat {
	if limit <= 0 {
		limit = 30
	}
	list := make([]*entity.DailyEmailAccountStat, 0, limit)
	_ = g.Model(string(entity.TbDailyEmailAccountStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyEmailAccountStats(list)
	return list
}

func reverseDailyEmailAccountStats(list []*entity.DailyEmailAccountStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyEmailAccountStat(data *entity.DailyEmailAccountStat) {
	if data == nil || data.ID == "" || dailyEmailAccountStatCacheMgr == nil {
		return
	}
	dailyEmailAccountStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}