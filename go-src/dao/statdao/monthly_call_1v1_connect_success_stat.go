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
	monthlyCall1v1ConnectSuccessStatCacheMgr     *cache.RowCache[*entity.MonthlyCall1v1ConnectSuccessStat]
	monthlyCall1v1ConnectSuccessStatListCacheMgr *cache.ListCache[*entity.MonthlyCall1v1ConnectSuccessStat]
)

func initMonthlyCall1v1ConnectSuccessStatDao() {
	monthlyCall1v1ConnectSuccessStatCacheMgr = cache.NewRowCache[*entity.MonthlyCall1v1ConnectSuccessStat]()
	monthlyCall1v1ConnectSuccessStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyCall1v1ConnectSuccessStat]()
}

func GetMonthlyCall1v1ConnectSuccessStatByMonth(month string) *entity.MonthlyCall1v1ConnectSuccessStat {
	return monthlyCall1v1ConnectSuccessStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyCall1v1ConnectSuccessStat, error) {
		var data *entity.MonthlyCall1v1ConnectSuccessStat
		_ = g.Model(string(entity.TbMonthlyCall1v1ConnectSuccessStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyCall1v1ConnectSuccessStat(month), nil
	})
}

func ListRecentMonthlyCall1v1ConnectSuccessStats(limit int) []*entity.MonthlyCall1v1ConnectSuccessStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyCall1v1ConnectSuccessStatKey(time.Now()))
	return monthlyCall1v1ConnectSuccessStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyCall1v1ConnectSuccessStat, error) {
		return loadRecentMonthlyCall1v1ConnectSuccessStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyCall1v1ConnectSuccessStatsFromDB(limit int) []*entity.MonthlyCall1v1ConnectSuccessStat {
	list := make([]*entity.MonthlyCall1v1ConnectSuccessStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyCall1v1ConnectSuccessStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyCall1v1ConnectSuccessStats(list)
	return list
}

func reverseMonthlyCall1v1ConnectSuccessStats(list []*entity.MonthlyCall1v1ConnectSuccessStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyCall1v1ConnectSuccessStat(data *entity.MonthlyCall1v1ConnectSuccessStat) {
	if data == nil || data.ID == "" || monthlyCall1v1ConnectSuccessStatCacheMgr == nil {
		return
	}
	monthlyCall1v1ConnectSuccessStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
