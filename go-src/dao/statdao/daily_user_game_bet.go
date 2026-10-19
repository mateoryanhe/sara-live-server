package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserGameBetCacheMgr *cache.RowCache[*entity.DailyUserGameBet]

func initDailyUserGameBetDao() {
	dailyUserGameBetCacheMgr = cache.NewRowCache[*entity.DailyUserGameBet]()
}

// TryRecordDailyGameBet 记录用户当日首次游戏下注;已下注过返回 false
func TryRecordDailyGameBet(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserGameBetId(date, userId)
	v := dailyUserGameBetCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserGameBet, error) {
		var row *entity.DailyUserGameBet
		_ = g.Model(string(entity.TbDailyUserGameBet)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserGameBet(date, userId), nil
		}
		return row, nil
	})
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishDailyUserGameBet(v)
	return true
}

func publishDailyUserGameBet(data *entity.DailyUserGameBet) {
	if data == nil || data.ID == "" || dailyUserGameBetCacheMgr == nil {
		return
	}
	dailyUserGameBetCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
