package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserGoldConsumeCacheMgr *cache.RowCache[*entity.WeeklyUserGoldConsume]

func initWeeklyUserGoldConsumeDao() {
	weeklyUserGoldConsumeCacheMgr = cache.NewRowCache[*entity.WeeklyUserGoldConsume]()
}

// TryRecordWeeklyGoldConsume 记录用户当周首次金币消费;已消费过返回 false
func TryRecordWeeklyGoldConsume(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserGoldConsumeId(week, userId)
	v := weeklyUserGoldConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserGoldConsume, error) {
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
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishWeeklyUserGoldConsume(v)
	return true
}
