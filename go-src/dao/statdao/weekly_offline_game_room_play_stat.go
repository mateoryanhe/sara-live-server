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
	weeklyOfflineGameRoomPlayStatCacheMgr     *cache.RowCache[*entity.WeeklyOfflineGameRoomPlayStat]
)

func initWeeklyOfflineGameRoomPlayStatDao() {
	weeklyOfflineGameRoomPlayStatCacheMgr = cache.NewRowCache[*entity.WeeklyOfflineGameRoomPlayStat]()
}

func GetWeeklyOfflineGameRoomPlayStatByWeek(week string) *entity.WeeklyOfflineGameRoomPlayStat {
	return weeklyOfflineGameRoomPlayStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyOfflineGameRoomPlayStat, error) {
		var data *entity.WeeklyOfflineGameRoomPlayStat
		_ = g.Model(string(entity.TbWeeklyOfflineGameRoomPlayStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyOfflineGameRoomPlayStat(week), nil
	})
}

func ListRecentWeeklyOfflineGameRoomPlayStats(limit int) []*entity.WeeklyOfflineGameRoomPlayStat {
	if limit <= 0 {
		limit = 12
	}
	return loadRecentWeeklyOfflineGameRoomPlayStatsFromDB(limit)
}

func loadRecentWeeklyOfflineGameRoomPlayStatsFromDB(limit int) []*entity.WeeklyOfflineGameRoomPlayStat {
	list := make([]*entity.WeeklyOfflineGameRoomPlayStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyOfflineGameRoomPlayStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyOfflineGameRoomPlayStats(list)
	return list
}

func reverseWeeklyOfflineGameRoomPlayStats(list []*entity.WeeklyOfflineGameRoomPlayStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyOfflineGameRoomPlayStat(data *entity.WeeklyOfflineGameRoomPlayStat) {
	if data == nil || data.ID == "" || weeklyOfflineGameRoomPlayStatCacheMgr == nil {
		return
	}
	weeklyOfflineGameRoomPlayStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
