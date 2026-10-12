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
	weeklyVoiceRoomMicRequestStatCacheMgr     *cache.RowCache[*entity.WeeklyVoiceRoomMicRequestStat]
	weeklyVoiceRoomMicRequestStatListCacheMgr *cache.ListCache[*entity.WeeklyVoiceRoomMicRequestStat]
)

func initWeeklyVoiceRoomMicRequestStatDao() {
	weeklyVoiceRoomMicRequestStatCacheMgr = cache.NewRowCache[*entity.WeeklyVoiceRoomMicRequestStat]()
	weeklyVoiceRoomMicRequestStatListCacheMgr = cache.NewPermanentListCache[*entity.WeeklyVoiceRoomMicRequestStat]()
}

func GetWeeklyVoiceRoomMicRequestStatByWeek(week string) *entity.WeeklyVoiceRoomMicRequestStat {
	return weeklyVoiceRoomMicRequestStatCacheMgr.MustGetRow(gctx.New(), week, func(ctx context.Context) (*entity.WeeklyVoiceRoomMicRequestStat, error) {
		var data *entity.WeeklyVoiceRoomMicRequestStat
		_ = g.Model(string(entity.TbWeeklyVoiceRoomMicRequestStat)).Unscoped().Where(g.Map{
			string(db.IdName): week,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewWeeklyVoiceRoomMicRequestStat(week), nil
	})
}

func ListRecentWeeklyVoiceRoomMicRequestStats(limit int) []*entity.WeeklyVoiceRoomMicRequestStat {
	if limit <= 0 {
		limit = 12
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatWeeklyVoiceRoomMicRequestStatKey(time.Now()))
	return weeklyVoiceRoomMicRequestStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.WeeklyVoiceRoomMicRequestStat, error) {
		return loadRecentWeeklyVoiceRoomMicRequestStatsFromDB(limit), nil
	})
}

func loadRecentWeeklyVoiceRoomMicRequestStatsFromDB(limit int) []*entity.WeeklyVoiceRoomMicRequestStat {
	list := make([]*entity.WeeklyVoiceRoomMicRequestStat, 0, limit)
	_ = g.Model(string(entity.TbWeeklyVoiceRoomMicRequestStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyVoiceRoomMicRequestStats(list)
	return list
}

func reverseWeeklyVoiceRoomMicRequestStats(list []*entity.WeeklyVoiceRoomMicRequestStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishWeeklyVoiceRoomMicRequestStat(data *entity.WeeklyVoiceRoomMicRequestStat) {
	if data == nil || data.ID == "" || weeklyVoiceRoomMicRequestStatCacheMgr == nil {
		return
	}
	weeklyVoiceRoomMicRequestStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
