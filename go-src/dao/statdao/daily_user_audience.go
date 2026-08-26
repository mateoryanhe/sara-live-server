package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserAudienceCacheMgr *cache.RowCache[*entity.DailyUserAudience]

func initDailyUserAudienceDao() {
	dailyUserAudienceCacheMgr = cache.NewRowCache[*entity.DailyUserAudience]()
}

// TryRecordDailyAudience 记录用户当日首次成为观众;已记录过返回 false
func TryRecordDailyAudience(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserAudienceId(date, userId)
	v := dailyUserAudienceCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserAudience, error) {
		var row *entity.DailyUserAudience
		_ = g.Model(string(entity.TbDailyUserAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserAudience(date, userId), nil
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
	publishDailyUserAudience(v)
	return true
}
