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
	weeklyLiveFirstFrameRenderStatCacheMgr     *cache.RowCache[*entity.WeeklyLiveFirstFrameRenderStat]
	weeklyLiveFirstFrameRenderStatListCacheMgr *cache.ListCache[*entity.WeeklyLiveFirstFrameRenderStat]
)

func initWeeklyLiveFirstFrameRenderStatDao() {
	weeklyLiveFirstFrameRenderStatCacheMgr = cache.NewRowCache[*entity.WeeklyLiveFirstFrameRenderStat]()
	weeklyLiveFirstFrameRenderStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyLiveFirstFrameRenderStat]()
}

func GetWeeklyLiveFirstFrameRenderStatByWeek(week string) *entity.WeeklyLiveFirstFrameRenderStat {
	return weeklyLiveFirstFrameRenderStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyLiveFirstFrameRenderStat, error) {
		var data *entity.WeeklyLiveFirstFrameRenderStat
		_ = g.Model(string(entity.TbWeeklyLiveFirstFrameRenderStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyLiveFirstFrameRenderStat(week), nil
	})
}

func ListRecentWeeklyLiveFirstFrameRenderStats(limit int) []*entity.WeeklyLiveFirstFrameRenderStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyLiveFirstFrameRenderStatKey(time.Now()))
	return weeklyLiveFirstFrameRenderStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyLiveFirstFrameRenderStat, error) {
		return loadRecentWeeklyLiveFirstFrameRenderStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyLiveFirstFrameRenderStatsFromDB(limit int) []*entity.WeeklyLiveFirstFrameRenderStat {
	list := make([]*entity.WeeklyLiveFirstFrameRenderStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyLiveFirstFrameRenderStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyLiveFirstFrameRenderStats(list)
	return list
}

func reverseWeeklyLiveFirstFrameRenderStats(list []*entity.WeeklyLiveFirstFrameRenderStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyLiveFirstFrameRenderStat(data *entity.WeeklyLiveFirstFrameRenderStat) {
	if data == nil || data.ID == "" || weeklyLiveFirstFrameRenderStatCacheMgr == nil {
		return
	}
	weeklyLiveFirstFrameRenderStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
