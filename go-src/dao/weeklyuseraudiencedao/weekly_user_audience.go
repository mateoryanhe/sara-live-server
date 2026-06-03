package weeklyuseraudiencedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var cacheMgr *cache.CacheMgr

func InitWeeklyUserAudienceDao() {
	cacheMgr = cache.NewCacheMgr()
}

func TryRecordAudience(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserAudienceId(week, userId)
	v := cacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
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
