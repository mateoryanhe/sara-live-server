package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserGameBetCacheMgr *cache.RowCache[*entity.MonthlyUserGameBet]

func initMonthlyUserGameBetDao() {
	monthlyUserGameBetCacheMgr = cache.NewRowCache[*entity.MonthlyUserGameBet]()
}

func TryRecordMonthlyGameBet(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserGameBetId(month, userId)
	v := monthlyUserGameBetCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserGameBet, error) {
		var row *entity.MonthlyUserGameBet
		_ = g.Model(string(entity.TbMonthlyUserGameBet)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserGameBet(month, userId), nil
		}
		return row, nil
	})
	if v == nil || v.CreatedAt != nil {
		return false
	}
	now := time.Now()
	v.SetCreatedAt(&now)
	publishMonthlyUserGameBet(v)
	return true
}

func publishMonthlyUserGameBet(data *entity.MonthlyUserGameBet) {
	if data == nil || data.ID == "" || monthlyUserGameBetCacheMgr == nil {
		return
	}
	monthlyUserGameBetCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
