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
	weeklyShortVideoPaywallClickStatCacheMgr     *cache.RowCache[*entity.WeeklyShortVideoPaywallClickStat]
	weeklyShortVideoPaywallClickStatListCacheMgr *cache.ListCache[*entity.WeeklyShortVideoPaywallClickStat]
)

func initWeeklyShortVideoPaywallClickStatDao() {
	weeklyShortVideoPaywallClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyShortVideoPaywallClickStat]()
	weeklyShortVideoPaywallClickStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyShortVideoPaywallClickStat]()
}

func GetWeeklyShortVideoPaywallClickStatByWeek(week string) *entity.WeeklyShortVideoPaywallClickStat {
	return weeklyShortVideoPaywallClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyShortVideoPaywallClickStat, error) {
		var data *entity.WeeklyShortVideoPaywallClickStat
		_ = g.Model(string(entity.TbWeeklyShortVideoPaywallClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyShortVideoPaywallClickStat(week), nil
	})
}

func ListRecentWeeklyShortVideoPaywallClickStats(limit int) []*entity.WeeklyShortVideoPaywallClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyShortVideoPaywallClickStatKey(time.Now()))
	return weeklyShortVideoPaywallClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyShortVideoPaywallClickStat, error) {
		return loadRecentWeeklyShortVideoPaywallClickStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyShortVideoPaywallClickStatsFromDB(limit int) []*entity.WeeklyShortVideoPaywallClickStat {
	list := make([]*entity.WeeklyShortVideoPaywallClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyShortVideoPaywallClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyShortVideoPaywallClickStats(list)
	return list
}

func reverseWeeklyShortVideoPaywallClickStats(list []*entity.WeeklyShortVideoPaywallClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyShortVideoPaywallClickStat(data *entity.WeeklyShortVideoPaywallClickStat) {
	if data == nil || data.ID == "" || weeklyShortVideoPaywallClickStatCacheMgr == nil {
		return
	}
	weeklyShortVideoPaywallClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
