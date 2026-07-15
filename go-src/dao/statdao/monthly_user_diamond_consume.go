package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var monthlyUserDiamondConsumeCacheMgr *cache.CacheMgr

func initMonthlyUserDiamondConsumeDao() {
	monthlyUserDiamondConsumeCacheMgr = cache.NewCacheMgr()
}

// TryRecordMonthlyDiamondConsume 记录用户当月首次钻石消费;已消费过返回 false
func TryRecordMonthlyDiamondConsume(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserDiamondConsumeId(month, userId)
	v := monthlyUserDiamondConsumeCacheMgr.GetData(id, func(ctx context.Context) (interface{}, error) {
		var row *entity.MonthlyUserDiamondConsume
		_ = g.Model(string(entity.TbMonthlyUserDiamondConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserDiamondConsume(month, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return false
	}
	data, _ := v.(*entity.MonthlyUserDiamondConsume)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
