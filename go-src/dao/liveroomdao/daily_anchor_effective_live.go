package liveroomdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var dailyAnchorEffectiveLiveCacheMgr *cache.CacheMgr

func initDailyAnchorEffectiveLiveDao() {
	dailyAnchorEffectiveLiveCacheMgr = cache.NewCacheMgr()
}

// GetDailyAnchorEffectiveLive 按日期+直播间获取每日有效直播次数(缓存;无则新建缓冲对象)
func GetDailyAnchorEffectiveLive(date string, roomId uint64) *entity.DailyAnchorEffectiveLive {
	if date == "" || roomId == 0 || dailyAnchorEffectiveLiveCacheMgr == nil {
		return nil
	}
	id := entity.BuildDailyAnchorEffectiveLiveId(date, roomId)
	v := dailyAnchorEffectiveLiveCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.DailyAnchorEffectiveLive
		_ = g.Model(string(entity.TbDailyAnchorEffectiveLive)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyAnchorEffectiveLive(date, roomId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.DailyAnchorEffectiveLive)
	return row
}

// AddDailyEffectiveLiveCount 当日有效直播次数+1(下播且本场时长达标时调用)
func AddDailyEffectiveLiveCount(roomId uint64, at time.Time) {
	if roomId == 0 {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyAnchorEffectiveLive(date, roomId)
	if row == nil {
		return
	}
	row.AddEffectiveLiveCount(1)
}
