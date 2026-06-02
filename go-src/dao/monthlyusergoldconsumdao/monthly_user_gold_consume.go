package monthlyusergoldconsumdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var cacheMgr *cache.CacheMgr

func InitMonthlyUserGoldConsumeDao() {
	cacheMgr = cache.NewCacheMgr()
}

func TryRecordConsume(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserGoldConsumeId(month, userId)
	v := cacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.MonthlyUserGoldConsume
		_ = g.Model(string(entity.TbMonthlyUserGoldConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserGoldConsume(month, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.MonthlyUserGoldConsume)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
