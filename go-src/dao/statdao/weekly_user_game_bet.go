package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserGameBetCacheMgr *cache.RowCache[*entity.WeeklyUserGameBet]

func initWeeklyUserGameBetDao() {
	weeklyUserGameBetCacheMgr = cache.NewRowCache[*entity.WeeklyUserGameBet]()
}

func TryRecordWeeklyGameBet(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserGameBetId(week, userId)
	v := weeklyUserGameBetCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserGameBet, error) {
		var row *entity.WeeklyUserGameBet
		_ = g.Model(string(entity.TbWeeklyUserGameBet)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserGameBet(week, userId), nil
		}
		return row, nil
	})
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishWeeklyUserGameBet(v)
	return true
}

func publishWeeklyUserGameBet(data *entity.WeeklyUserGameBet) {
	if data == nil || data.ID == "" || weeklyUserGameBetCacheMgr == nil {
		return
	}
	weeklyUserGameBetCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
