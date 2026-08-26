package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserDiamondConsumeCacheMgr *cache.RowCache[*entity.MonthlyUserDiamondConsume]

func initMonthlyUserDiamondConsumeDao() {
	monthlyUserDiamondConsumeCacheMgr = cache.NewRowCache[*entity.MonthlyUserDiamondConsume]()
}

// TryRecordMonthlyDiamondConsume 记录用户当月首次钻石消费;已消费过返回 false
func TryRecordMonthlyDiamondConsume(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserDiamondConsumeId(month, userId)
	v := monthlyUserDiamondConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserDiamondConsume, error) {
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
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishMonthlyUserDiamondConsume(v)
	return true
}
