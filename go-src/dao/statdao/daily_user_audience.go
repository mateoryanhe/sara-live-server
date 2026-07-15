package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var dailyUserAudienceCacheMgr *cache.CacheMgr

func initDailyUserAudienceDao() {
	dailyUserAudienceCacheMgr = cache.NewCacheMgr()
}

// TryRecordDailyAudience 记录用户当日首次成为观众;已记录过返回 false
func TryRecordDailyAudience(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserAudienceId(date, userId)
	v := dailyUserAudienceCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.DailyUserAudience
		_ = g.Model(string(entity.TbDailyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserAudience(date, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.DailyUserAudience)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
