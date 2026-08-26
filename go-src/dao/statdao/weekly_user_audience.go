package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserAudienceCacheMgr *cache.RowCache[*entity.WeeklyUserAudience]

func initWeeklyUserAudienceDao() {
	weeklyUserAudienceCacheMgr = cache.NewRowCache[*entity.WeeklyUserAudience]()
}

// TryRecordWeeklyAudience 记录用户当周首次成为观众;已记录过返回 false
func TryRecordWeeklyAudience(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserAudienceId(week, userId)
	v := weeklyUserAudienceCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserAudience, error) {
		var row *entity.WeeklyUserAudience
		_ = g.Model(string(entity.TbWeeklyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserAudience(week, userId), nil
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
	publishWeeklyUserAudience(v)
	return true
}
