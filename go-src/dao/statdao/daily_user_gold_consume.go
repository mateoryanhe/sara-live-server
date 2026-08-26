package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserGoldConsumeCacheMgr *cache.RowCache[*entity.DailyUserGoldConsume]

func initDailyUserGoldConsumeDao() {
	dailyUserGoldConsumeCacheMgr = cache.NewRowCache[*entity.DailyUserGoldConsume]()
}

// TryRecordDailyGoldConsume 记录用户当日首次金币消费;已消费过返回 false
func TryRecordDailyGoldConsume(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserGoldConsumeId(date, userId)
	v := dailyUserGoldConsumeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserGoldConsume, error) {
		var row *entity.DailyUserGoldConsume
		_ = g.Model(string(entity.TbDailyUserGoldConsume)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserGoldConsume(date, userId), nil
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
	publishDailyUserGoldConsume(v)
	return true
}
