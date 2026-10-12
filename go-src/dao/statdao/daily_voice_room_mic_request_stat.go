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
	dailyVoiceRoomMicRequestStatCacheMgr     *cache.RowCache[*entity.DailyVoiceRoomMicRequestStat]
	dailyVoiceRoomMicRequestStatListCacheMgr *cache.ListCache[*entity.DailyVoiceRoomMicRequestStat]
)

func initDailyVoiceRoomMicRequestStatDao() {
	dailyVoiceRoomMicRequestStatCacheMgr = cache.NewRowCache[*entity.DailyVoiceRoomMicRequestStat]()
	dailyVoiceRoomMicRequestStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyVoiceRoomMicRequestStat]()
}

func GetDailyVoiceRoomMicRequestStatByDate(date string) *entity.DailyVoiceRoomMicRequestStat {
	return dailyVoiceRoomMicRequestStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyVoiceRoomMicRequestStat, error) {
		var data *entity.DailyVoiceRoomMicRequestStat
		_ = g.Model(string(entity.TbDailyVoiceRoomMicRequestStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyVoiceRoomMicRequestStat(date), nil
	})
}

func ListRecentDailyVoiceRoomMicRequestStats(limit int) []*entity.DailyVoiceRoomMicRequestStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyVoiceRoomMicRequestStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyVoiceRoomMicRequestStat, error) {
		return loadRecentDailyVoiceRoomMicRequestStatsFromDB(limit), nil
	})
}

func loadRecentDailyVoiceRoomMicRequestStatsFromDB(limit int) []*entity.DailyVoiceRoomMicRequestStat {
	list := make([]*entity.DailyVoiceRoomMicRequestStat, 0, limit)
	_ = g.Model(string(entity.TbDailyVoiceRoomMicRequestStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyVoiceRoomMicRequestStats(list)
	return list
}

func reverseDailyVoiceRoomMicRequestStats(list []*entity.DailyVoiceRoomMicRequestStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyVoiceRoomMicRequestStat(data *entity.DailyVoiceRoomMicRequestStat) {
	if data == nil || data.ID == "" || dailyVoiceRoomMicRequestStatCacheMgr == nil {
		return
	}
	dailyVoiceRoomMicRequestStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
