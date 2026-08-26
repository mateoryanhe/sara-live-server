package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserDiamondConsumeCacheMgr *cache.RowCache[*entity.WeeklyUserDiamondConsume]

func initWeeklyUserDiamondConsumeDao() {
	weeklyUserDiamondConsumeCacheMgr = cache.NewRowCache[*entity.WeeklyUserDiamondConsume]()
}

// TryRecordWeeklyDiamondConsume 记录用户当周首次钻石消费;已消费过返回 false
func TryRecordWeeklyDiamondConsume(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserDiamondConsumeId(week, userId)
	v := weeklyUserDiamondConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserDiamondConsume, error) {
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
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishWeeklyUserDiamondConsume(v)
	return true
}
