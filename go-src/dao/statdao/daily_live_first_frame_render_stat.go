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
	dailyLiveFirstFrameRenderStatCacheMgr     *cache.RowCache[*entity.DailyLiveFirstFrameRenderStat]
	dailyLiveFirstFrameRenderStatListCacheMgr *cache.ListCache[*entity.DailyLiveFirstFrameRenderStat]
)

func initDailyLiveFirstFrameRenderStatDao() {
	dailyLiveFirstFrameRenderStatCacheMgr = cache.NewRowCache[*entity.DailyLiveFirstFrameRenderStat]()
	dailyLiveFirstFrameRenderStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyLiveFirstFrameRenderStat]()
}

func GetDailyLiveFirstFrameRenderStatByDate(date string) *entity.DailyLiveFirstFrameRenderStat {
	return dailyLiveFirstFrameRenderStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyLiveFirstFrameRenderStat, error) {
		var data *entity.DailyLiveFirstFrameRenderStat
		_ = g.Model(string(entity.TbDailyLiveFirstFrameRenderStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyLiveFirstFrameRenderStat(date), nil
	})
}

func ListRecentDailyLiveFirstFrameRenderStats(limit int) []*entity.DailyLiveFirstFrameRenderStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyLiveFirstFrameRenderStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyLiveFirstFrameRenderStat, error) {
		return loadRecentDailyLiveFirstFrameRenderStatsFromDB(limit), nil
	})
}

func loadRecentDailyLiveFirstFrameRenderStatsFromDB(limit int) []*entity.DailyLiveFirstFrameRenderStat {
	list := make([]*entity.DailyLiveFirstFrameRenderStat, 0, limit)
	_ = g.Model(string(entity.TbDailyLiveFirstFrameRenderStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyLiveFirstFrameRenderStats(list)
	return list
}

func reverseDailyLiveFirstFrameRenderStats(list []*entity.DailyLiveFirstFrameRenderStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyLiveFirstFrameRenderStat(data *entity.DailyLiveFirstFrameRenderStat) {
	if data == nil || data.ID == "" || dailyLiveFirstFrameRenderStatCacheMgr == nil {
		return
	}
	dailyLiveFirstFrameRenderStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
