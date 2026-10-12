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
	monthlyVoiceRoomEnterStatCacheMgr     *cache.RowCache[*entity.MonthlyVoiceRoomEnterStat]
	monthlyVoiceRoomEnterStatListCacheMgr *cache.ListCache[*entity.MonthlyVoiceRoomEnterStat]
)

func initMonthlyVoiceRoomEnterStatDao() {
	monthlyVoiceRoomEnterStatCacheMgr = cache.NewRowCache[*entity.MonthlyVoiceRoomEnterStat]()
	monthlyVoiceRoomEnterStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyVoiceRoomEnterStat]()
}

func GetMonthlyVoiceRoomEnterStatByMonth(month string) *entity.MonthlyVoiceRoomEnterStat {
	return monthlyVoiceRoomEnterStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyVoiceRoomEnterStat, error) {
		var data *entity.MonthlyVoiceRoomEnterStat
		_ = g.Model(string(entity.TbMonthlyVoiceRoomEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyVoiceRoomEnterStat(month), nil
	})
}

func ListRecentMonthlyVoiceRoomEnterStats(limit int) []*entity.MonthlyVoiceRoomEnterStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyVoiceRoomEnterStatKey(time.Now()))
	return monthlyVoiceRoomEnterStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyVoiceRoomEnterStat, error) {
		return loadRecentMonthlyVoiceRoomEnterStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyVoiceRoomEnterStatsFromDB(limit int) []*entity.MonthlyVoiceRoomEnterStat {
	list := make([]*entity.MonthlyVoiceRoomEnterStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyVoiceRoomEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyVoiceRoomEnterStats(list)
	return list
}

func reverseMonthlyVoiceRoomEnterStats(list []*entity.MonthlyVoiceRoomEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyVoiceRoomEnterStat(data *entity.MonthlyVoiceRoomEnterStat) {
	if data == nil || data.ID == "" || monthlyVoiceRoomEnterStatCacheMgr == nil {
		return
	}
	monthlyVoiceRoomEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
