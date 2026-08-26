package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserAudienceCacheMgr *cache.RowCache[*entity.MonthlyUserAudience]

func initMonthlyUserAudienceDao() {
	monthlyUserAudienceCacheMgr = cache.NewRowCache[*entity.MonthlyUserAudience]()
}

// TryRecordMonthlyAudience 记录用户当月首次成为观众;已记录过返回 false
func TryRecordMonthlyAudience(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserAudienceId(month, userId)
	v := monthlyUserAudienceCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserAudience, error) {
		var row *entity.MonthlyUserAudience
		_ = g.Model(string(entity.TbMonthlyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserAudience(month, userId), nil
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
	publishMonthlyUserAudience(v)
	return true
}
