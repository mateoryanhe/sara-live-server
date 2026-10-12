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
	dailyVoiceRoomEnterStatCacheMgr     *cache.RowCache[*entity.DailyVoiceRoomEnterStat]
	dailyVoiceRoomEnterStatListCacheMgr *cache.ListCache[*entity.DailyVoiceRoomEnterStat]
)

func initDailyVoiceRoomEnterStatDao() {
	dailyVoiceRoomEnterStatCacheMgr = cache.NewRowCache[*entity.DailyVoiceRoomEnterStat]()
	dailyVoiceRoomEnterStatListCacheMgr = cache.NewPermanentListCache[*entity.DailyVoiceRoomEnterStat]()
}

func GetDailyVoiceRoomEnterStatByDate(date string) *entity.DailyVoiceRoomEnterStat {
	return dailyVoiceRoomEnterStatCacheMgr.MustGetRow(gctx.New(), date, func(ctx context.Context) (*entity.DailyVoiceRoomEnterStat, error) {
		var data *entity.DailyVoiceRoomEnterStat
		_ = g.Model(string(entity.TbDailyVoiceRoomEnterStat)).Unscoped().Where(g.Map{
			string(db.IdName): date,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewDailyVoiceRoomEnterStat(date), nil
	})
}

func ListRecentDailyVoiceRoomEnterStats(limit int) []*entity.DailyVoiceRoomEnterStat {
	if limit <= 0 {
		limit = 30
	}
	key := fmt.Sprintf("%d:%s", limit, entity.FormatDailyLoginStatDate(time.Now()))
	return dailyVoiceRoomEnterStatListCacheMgr.MustGetList(gctx.New(), key, func(ctx context.Context) ([]*entity.DailyVoiceRoomEnterStat, error) {
		return loadRecentDailyVoiceRoomEnterStatsFromDB(limit), nil
	})
}

func loadRecentDailyVoiceRoomEnterStatsFromDB(limit int) []*entity.DailyVoiceRoomEnterStat {
	list := make([]*entity.DailyVoiceRoomEnterStat, 0, limit)
	_ = g.Model(string(entity.TbDailyVoiceRoomEnterStat)).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyVoiceRoomEnterStats(list)
	return list
}

func reverseDailyVoiceRoomEnterStats(list []*entity.DailyVoiceRoomEnterStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func PublishDailyVoiceRoomEnterStat(data *entity.DailyVoiceRoomEnterStat) {
	if data == nil || data.ID == "" || dailyVoiceRoomEnterStatCacheMgr == nil {
		return
	}
	dailyVoiceRoomEnterStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
