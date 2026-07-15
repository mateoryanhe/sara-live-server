package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var weeklyUserDiamondConsumeCacheMgr *cache.CacheMgr

func initWeeklyUserDiamondConsumeDao() {
	weeklyUserDiamondConsumeCacheMgr = cache.NewCacheMgr()
}

// TryRecordWeeklyDiamondConsume 记录用户当周首次钻石消费;已消费过返回 false
func TryRecordWeeklyDiamondConsume(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserDiamondConsumeId(week, userId)
	v := weeklyUserDiamondConsumeCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.WeeklyUserDiamondConsume
		_ = g.Model(string(entity.TbWeeklyUserDiamondConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserDiamondConsume(week, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.WeeklyUserDiamondConsume)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
