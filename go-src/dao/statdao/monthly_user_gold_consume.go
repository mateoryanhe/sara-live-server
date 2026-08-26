package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserGoldConsumeCacheMgr *cache.RowCache[*entity.MonthlyUserGoldConsume]

func initMonthlyUserGoldConsumeDao() {
	monthlyUserGoldConsumeCacheMgr = cache.NewRowCache[*entity.MonthlyUserGoldConsume]()
}

// TryRecordMonthlyGoldConsume 记录用户当月首次金币消费;已消费过返回 false
func TryRecordMonthlyGoldConsume(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserGoldConsumeId(month, userId)
	v := monthlyUserGoldConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserGoldConsume, error) {
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
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishMonthlyUserGoldConsume(v)
	return true
}
