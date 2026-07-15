package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var dailyUserGoldConsumeCacheMgr *cache.CacheMgr

func initDailyUserGoldConsumeDao() {
	dailyUserGoldConsumeCacheMgr = cache.NewCacheMgr()
}

// TryRecordDailyGoldConsume 记录用户当日首次金币消费;已消费过返回 false
func TryRecordDailyGoldConsume(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserGoldConsumeId(date, userId)
	v := dailyUserGoldConsumeCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
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
