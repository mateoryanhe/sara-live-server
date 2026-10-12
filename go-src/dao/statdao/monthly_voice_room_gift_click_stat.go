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
	monthlyVoiceRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.MonthlyVoiceRoomGiftClickStat]
	monthlyVoiceRoomGiftClickStatListCacheMgr *cache.ListCache[*entity.MonthlyVoiceRoomGiftClickStat]
)

func initMonthlyVoiceRoomGiftClickStatDao() {
	monthlyVoiceRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.MonthlyVoiceRoomGiftClickStat]()
	monthlyVoiceRoomGiftClickStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyVoiceRoomGiftClickStat]()
}

func GetMonthlyVoiceRoomGiftClickStatByMonth(month string) *entity.MonthlyVoiceRoomGiftClickStat {
	return monthlyVoiceRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyVoiceRoomGiftClickStat, error) {
		var data *entity.MonthlyVoiceRoomGiftClickStat
		_ = g.Model(string(entity.TbMonthlyVoiceRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyVoiceRoomGiftClickStat(month), nil
	})
}

func ListRecentMonthlyVoiceRoomGiftClickStats(limit int) []*entity.MonthlyVoiceRoomGiftClickStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyVoiceRoomGiftClickStatKey(time.Now()))
	return monthlyVoiceRoomGiftClickStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyVoiceRoomGiftClickStat, error) {
		return loadRecentMonthlyVoiceRoomGiftClickStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyVoiceRoomGiftClickStatsFromDB(limit int) []*entity.MonthlyVoiceRoomGiftClickStat {
	list := make([]*entity.MonthlyVoiceRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyVoiceRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyVoiceRoomGiftClickStats(list)
	return list
}

func reverseMonthlyVoiceRoomGiftClickStats(list []*entity.MonthlyVoiceRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyVoiceRoomGiftClickStat(data *entity.MonthlyVoiceRoomGiftClickStat) {
	if data == nil || data.ID == "" || monthlyVoiceRoomGiftClickStatCacheMgr == nil {
		return
	}
	monthlyVoiceRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
