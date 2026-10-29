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
	weeklyVoiceRoomGiftClickStatCacheMgr     *cache.RowCache[*entity.WeeklyVoiceRoomGiftClickStat]
)

func initWeeklyVoiceRoomGiftClickStatDao() {
	weeklyVoiceRoomGiftClickStatCacheMgr = cache.NewRowCache[*entity.WeeklyVoiceRoomGiftClickStat]()
}

func GetWeeklyVoiceRoomGiftClickStatByWeek(week string) *entity.WeeklyVoiceRoomGiftClickStat {
	return weeklyVoiceRoomGiftClickStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyVoiceRoomGiftClickStat, error) {
		var data *entity.WeeklyVoiceRoomGiftClickStat
		_ = g.Model(string(entity.TbWeeklyVoiceRoomGiftClickStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyVoiceRoomGiftClickStat(week), nil
	})
}

func ListRecentWeeklyVoiceRoomGiftClickStats(limit int) []*entity.WeeklyVoiceRoomGiftClickStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyVoiceRoomGiftClickStatsFromDB(limit)
}

func loadRecentWeeklyVoiceRoomGiftClickStatsFromDB(limit int) []*entity.WeeklyVoiceRoomGiftClickStat {
	list := make([]*entity.WeeklyVoiceRoomGiftClickStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyVoiceRoomGiftClickStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyVoiceRoomGiftClickStats(list)
	return list
}

func reverseWeeklyVoiceRoomGiftClickStats(list []*entity.WeeklyVoiceRoomGiftClickStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyVoiceRoomGiftClickStat(data *entity.WeeklyVoiceRoomGiftClickStat) {
	if data == nil || data.ID == "" || weeklyVoiceRoomGiftClickStatCacheMgr == nil {
		return
	}
	weeklyVoiceRoomGiftClickStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
