package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var (
	dailyVoiceRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.DailyVoiceRoomGiftClickStat]
)

func initDailyVoiceRoomGiftClickStatDao() {
	dailyVoiceRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.DailyVoiceRoomGiftClickStat]()
}

func GetDailyVoiceRoomGiftClickStatByDate(date string) *entity.DailyVoiceRoomGiftClickStat {
	return dailyVoiceRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyVoiceRoomGiftClickStat, error) {
		var data *entity.DailyVoiceRoomGiftClickStat
		_ = g.Model(string(entity.TbDailyVoiceRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyVoiceRoomGiftClickStat(date), nil
	})
}

func ListRecentDailyVoiceRoomGiftClickStats(limit int) []*entity.DailyVoiceRoomGiftClickStat {
	if limit <= 0 {
		limit = 30
	}
	// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。
	return loadRecentDailyVoiceRoomGiftClickStatsFromDB(limit)
}

func loadRecentDailyVoiceRoomGiftClickStatsFromDB(limit int) []*entity.DailyVoiceRoomGiftClickStat {
	list := make([]*entity.DailyVoiceRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbDailyVoiceRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyVoiceRoomGiftClickStats(list)
	return list
}

func reverseDailyVoiceRoomGiftClickStats(list []*entity.DailyVoiceRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyVoiceRoomGiftClickStat(data *entity.DailyVoiceRoomGiftClickStat) {
	if data == nil || data.ID == "" || dailyVoiceRoomGiftClickStatCacheMgr == nil {
		return
	}
	dailyVoiceRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
