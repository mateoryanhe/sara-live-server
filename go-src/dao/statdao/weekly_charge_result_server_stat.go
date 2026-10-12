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
	weeklyChargeResultServerStatCacheMgr     *cache.RowCache[*entity.WeeklyChargeResultServerStat]
	weeklyChargeResultServerStatListCacheMgr *cache.ListCache[*entity.WeeklyChargeResultServerStat]
)

func initWeeklyChargeResultServerStatDao() {
	weeklyChargeResultServerStatCacheMgr = cache.NewRowCache[*entity.WeeklyChargeResultServerStat]()
	weeklyChargeResultServerStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyChargeResultServerStat]()
}

func GetWeeklyChargeResultServerStatByWeek(week string) *entity.WeeklyChargeResultServerStat {
	return weeklyChargeResultServerStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyChargeResultServerStat, error) {
		var data *entity.WeeklyChargeResultServerStat
		_ = g.Model(string(entity.TbWeeklyChargeResultServerStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyChargeResultServerStat(week), nil
	})
}

func ListRecentWeeklyChargeResultServerStats(limit int) []*entity.WeeklyChargeResultServerStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyChargeResultServerStatKey(time.Now()))
	return weeklyChargeResultServerStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyChargeResultServerStat, error) {
		return loadRecentWeeklyChargeResultServerStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyChargeResultServerStatsFromDB(limit int) []*entity.WeeklyChargeResultServerStat {
	list := make([]*entity.WeeklyChargeResultServerStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyChargeResultServerStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyChargeResultServerStats(list)
	return list
}

func reverseWeeklyChargeResultServerStats(list []*entity.WeeklyChargeResultServerStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyChargeResultServerStat(data *entity.WeeklyChargeResultServerStat) {
	if data == nil || data.ID == "" || weeklyChargeResultServerStatCacheMgr == nil {
		return
	}
	weeklyChargeResultServerStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
