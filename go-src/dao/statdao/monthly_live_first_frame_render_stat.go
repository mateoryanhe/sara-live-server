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
	monthlyLiveFirstFrameRenderStatCacheMgr     *cache.RowCache[*entity.MonthlyLiveFirstFrameRenderStat]
	monthlyLiveFirstFrameRenderStatListCacheMgr *cache.ListCache[*entity.MonthlyLiveFirstFrameRenderStat]
)

func initMonthlyLiveFirstFrameRenderStatDao() {
	monthlyLiveFirstFrameRenderStatCacheMgr = cache.NewRowCache[*entity.MonthlyLiveFirstFrameRenderStat]()
	monthlyLiveFirstFrameRenderStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyLiveFirstFrameRenderStat]()
}

func GetMonthlyLiveFirstFrameRenderStatByMonth(month string) *entity.MonthlyLiveFirstFrameRenderStat {
	return monthlyLiveFirstFrameRenderStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyLiveFirstFrameRenderStat, error) {
		var data *entity.MonthlyLiveFirstFrameRenderStat
		_ = g.Model(string(entity.TbMonthlyLiveFirstFrameRenderStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyLiveFirstFrameRenderStat(month), nil
	})
}

func ListRecentMonthlyLiveFirstFrameRenderStats(limit int) []*entity.MonthlyLiveFirstFrameRenderStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyLiveFirstFrameRenderStatKey(time.Now()))
	return monthlyLiveFirstFrameRenderStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyLiveFirstFrameRenderStat, error) {
		return loadRecentMonthlyLiveFirstFrameRenderStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyLiveFirstFrameRenderStatsFromDB(limit int) []*entity.MonthlyLiveFirstFrameRenderStat {
	list := make([]*entity.MonthlyLiveFirstFrameRenderStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyLiveFirstFrameRenderStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyLiveFirstFrameRenderStats(list)
	return list
}

func reverseMonthlyLiveFirstFrameRenderStats(list []*entity.MonthlyLiveFirstFrameRenderStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyLiveFirstFrameRenderStat(data *entity.MonthlyLiveFirstFrameRenderStat) {
	if data == nil || data.ID == "" || monthlyLiveFirstFrameRenderStatCacheMgr == nil {
		return
	}
	monthlyLiveFirstFrameRenderStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
