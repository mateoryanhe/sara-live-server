package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserRechargeCacheMgr *cache.CacheMgr

func initMonthlyUserRechargeDao() {
	monthlyUserRechargeCacheMgr = cache.NewCacheMgr()
}

func getMonthlyUserRechargeById(id string, month string, userId uint64) *entity.MonthlyUserRecharge {
	v := monthlyUserRechargeCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.MonthlyUserRecharge
		_ = g.Model(string(entity.TbMonthlyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserRecharge(month, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.MonthlyUserRecharge)
	return row
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
	return true
}
