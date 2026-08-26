package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserDiamondConsumeCacheMgr *cache.RowCache[*entity.DailyUserDiamondConsume]

func initDailyUserDiamondConsumeDao() {
	dailyUserDiamondConsumeCacheMgr = cache.NewRowCache[*entity.DailyUserDiamondConsume]()
}

// TryRecordDailyDiamondConsume 记录用户当日首次钻石消费;已消费过返回 false
func TryRecordDailyDiamondConsume(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserDiamondConsumeId(date, userId)
	v := dailyUserDiamondConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserDiamondConsume, error) {
		var row *entity.DailyUserDiamondConsume
		_ = g.Model(string(entity.TbDailyUserDiamondConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserDiamondConsume(date, userId), nil
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
	publishDailyUserDiamondConsume(v)
	return true
}
