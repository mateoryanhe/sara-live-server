package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var weeklyUserAudienceCacheMgr *cache.CacheMgr

func initWeeklyUserAudienceDao() {
	weeklyUserAudienceCacheMgr = cache.NewCacheMgr()
}

// TryRecordWeeklyAudience 记录用户当周首次成为观众;已记录过返回 false
func TryRecordWeeklyAudience(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserAudienceId(week, userId)
	v := weeklyUserAudienceCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.WeeklyUserAudience
		_ = g.Model(string(entity.TbWeeklyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserAudience(week, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.WeeklyUserAudience)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
