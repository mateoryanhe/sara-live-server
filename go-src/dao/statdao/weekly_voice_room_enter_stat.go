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
	weeklyVoiceRoomEnterStatCacheMgr     *cache.RowCache[*entity.WeeklyVoiceRoomEnterStat]
)

func initWeeklyVoiceRoomEnterStatDao() {
	weeklyVoiceRoomEnterStatCacheMgr = cache.NewRowCache[*entity.WeeklyVoiceRoomEnterStat]()
}

func GetWeeklyVoiceRoomEnterStatByWeek(week string) *entity.WeeklyVoiceRoomEnterStat {
	return weeklyVoiceRoomEnterStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyVoiceRoomEnterStat, error) {
		var data *entity.WeeklyVoiceRoomEnterStat
		_ = g.Model(string(entity.TbWeeklyVoiceRoomEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyVoiceRoomEnterStat(week), nil
	})
}

func ListRecentWeeklyVoiceRoomEnterStats(limit int) []*entity.WeeklyVoiceRoomEnterStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyVoiceRoomEnterStatsFromDB(limit)
}

func loadRecentWeeklyVoiceRoomEnterStatsFromDB(limit int) []*entity.WeeklyVoiceRoomEnterStat {
	list := make([]*entity.WeeklyVoiceRoomEnterStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyVoiceRoomEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyVoiceRoomEnterStats(list)
	return list
}

func reverseWeeklyVoiceRoomEnterStats(list []*entity.WeeklyVoiceRoomEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyVoiceRoomEnterStat(data *entity.WeeklyVoiceRoomEnterStat) {
	if data == nil || data.ID == "" || weeklyVoiceRoomEnterStatCacheMgr == nil {
		return
	}
	weeklyVoiceRoomEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
