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
	weeklyGameLiveRoomJoinStatCacheMgr     *cache.RowCache[*entity.WeeklyGameLiveRoomJoinStat]
	weeklyGameLiveRoomJoinStatListCacheMgr *cache.ListCache[*entity.WeeklyGameLiveRoomJoinStat]
)

func initWeeklyGameLiveRoomJoinStatDao() {
	weeklyGameLiveRoomJoinStatCacheMgr = cache.NewRowCache[*entity.WeeklyGameLiveRoomJoinStat]()
	weeklyGameLiveRoomJoinStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyGameLiveRoomJoinStat]()
}

func GetWeeklyGameLiveRoomJoinStatByWeek(week string) *entity.WeeklyGameLiveRoomJoinStat {
	return weeklyGameLiveRoomJoinStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyGameLiveRoomJoinStat, error) {
		var data *entity.WeeklyGameLiveRoomJoinStat
		_ = g.Model(string(entity.TbWeeklyGameLiveRoomJoinStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyGameLiveRoomJoinStat(week), nil
	})
}

func ListRecentWeeklyGameLiveRoomJoinStats(limit int) []*entity.WeeklyGameLiveRoomJoinStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyGameLiveRoomJoinStatKey(time.Now()))
	return weeklyGameLiveRoomJoinStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyGameLiveRoomJoinStat, error) {
		return loadRecentWeeklyGameLiveRoomJoinStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyGameLiveRoomJoinStatsFromDB(limit int) []*entity.WeeklyGameLiveRoomJoinStat {
	list := make([]*entity.WeeklyGameLiveRoomJoinStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyGameLiveRoomJoinStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyGameLiveRoomJoinStats(list)
	return list
}

func reverseWeeklyGameLiveRoomJoinStats(list []*entity.WeeklyGameLiveRoomJoinStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyGameLiveRoomJoinStat(data *entity.WeeklyGameLiveRoomJoinStat) {
	if data == nil || data.ID == "" || weeklyGameLiveRoomJoinStatCacheMgr == nil {
		return
	}
	weeklyGameLiveRoomJoinStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
