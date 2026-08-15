package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserGoldConsumeCacheMgr *cache.CacheMgr

func initWeeklyUserGoldConsumeDao() {
	weeklyUserGoldConsumeCacheMgr = cache.NewCacheMgr()
}

// TryRecordWeeklyGoldConsume 记录用户当周首次金币消费;已消费过返回 false
func TryRecordWeeklyGoldConsume(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserGoldConsumeId(week, userId)
	v := weeklyUserGoldConsumeCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.WeeklyUserGoldConsume
		_ = g.Model(string(entity.TbWeeklyUserGoldConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserGoldConsume(week, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.WeeklyUserGoldConsume)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
