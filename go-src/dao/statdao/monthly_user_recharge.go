package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserRechargeCacheMgr *cache.RowCache[*entity.MonthlyUserRecharge]

func initMonthlyUserRechargeDao() {
	monthlyUserRechargeCacheMgr = cache.NewRowCache[*entity.MonthlyUserRecharge]()
}

func getMonthlyUserRechargeById(id string, month string, userId uint64) *entity.MonthlyUserRecharge {
	v := monthlyUserRechargeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserRecharge, error) {
		var row *entity.MonthlyUserRecharge
		_ = g.Model(string(entity.TbMonthlyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserRecharge(month, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordMonthlyRecharge 记录用户当月首次充值;本月已充值过返回 false
func TryRecordMonthlyRecharge(month string, userId uint64) bool {
	if userId == 0 || month == "" {
		return false
	}
	id := entity.BuildMonthlyUserRechargeId(month, userId)
	data := getMonthlyUserRechargeById(id, month, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishMonthlyUserRecharge(data)
	return true
}
