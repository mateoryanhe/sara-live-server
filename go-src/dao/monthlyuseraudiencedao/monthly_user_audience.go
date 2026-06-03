package monthlyuseraudiencedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var cacheMgr *cache.CacheMgr

func InitMonthlyUserAudienceDao() {
	cacheMgr = cache.NewCacheMgr()
}

func TryRecordAudience(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserAudienceId(month, userId)
	v := cacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.MonthlyUserAudience
		_ = g.Model(string(entity.TbMonthlyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserAudience(month, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.MonthlyUserAudience)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
