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
	monthlyVoiceRoomMicRequestStatCacheMgr     *cache.RowCache[*entity.MonthlyVoiceRoomMicRequestStat]
	monthlyVoiceRoomMicRequestStatListCacheMgr *cache.ListCache[*entity.MonthlyVoiceRoomMicRequestStat]
)

func initMonthlyVoiceRoomMicRequestStatDao() {
	monthlyVoiceRoomMicRequestStatCacheMgr = cache.NewRowCache[*entity.MonthlyVoiceRoomMicRequestStat]()
	monthlyVoiceRoomMicRequestStatListCacheMgr = cache.NewPermanentListCache[*entity.MonthlyVoiceRoomMicRequestStat]()
}

func GetMonthlyVoiceRoomMicRequestStatByMonth(month string) *entity.MonthlyVoiceRoomMicRequestStat {
	return monthlyVoiceRoomMicRequestStatCacheMgr.MustGetRow(gctx.New(), month, func(ctx context.Context) (*entity.MonthlyVoiceRoomMicRequestStat, error) {
		var data *entity.MonthlyVoiceRoomMicRequestStat
		_ = g.Model(string(entity.TbMonthlyVoiceRoomMicRequestStat)).Unscoped().Where(g.Map{
			string(db.IdName): month,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMonthlyVoiceRoomMicRequestStat(month), nil
	})
}

func ListRecentMonthlyVoiceRoomMicRequestStats(limit int) []*entity.MonthlyVoiceRoomMicRequestStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatMonthlyVoiceRoomMicRequestStatKey(time.Now()))
	return monthlyVoiceRoomMicRequestStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.MonthlyVoiceRoomMicRequestStat, error) {
		return loadRecentMonthlyVoiceRoomMicRequestStatsFromDB(limit), nil
	})
}

func loadRecentMonthlyVoiceRoomMicRequestStatsFromDB(limit int) []*entity.MonthlyVoiceRoomMicRequestStat {
	list := make([]*entity.MonthlyVoiceRoomMicRequestStat, 0, limit)
	_ = g.Model(string(entity.TbMonthlyVoiceRoomMicRequestStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseMonthlyVoiceRoomMicRequestStats(list)
	return list
}

func reverseMonthlyVoiceRoomMicRequestStats(list []*entity.MonthlyVoiceRoomMicRequestStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishMonthlyVoiceRoomMicRequestStat(data *entity.MonthlyVoiceRoomMicRequestStat) {
	if data == nil || data.ID == "" || monthlyVoiceRoomMicRequestStatCacheMgr == nil {
		return
	}
	monthlyVoiceRoomMicRequestStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
