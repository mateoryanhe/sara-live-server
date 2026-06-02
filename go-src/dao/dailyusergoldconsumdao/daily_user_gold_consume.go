package dailyusergoldconsumdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var cacheMgr *cache.CacheMgr

func InitDailyUserGoldConsumeDao() {
	cacheMgr = cache.NewCacheMgr()
}

func TryRecordConsume(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserGoldConsumeId(date, userId)
	v := cacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.DailyUserGoldConsume
		_ = g.Model(string(entity.TbDailyUserGoldConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserGoldConsume(date, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.DailyUserGoldConsume)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
