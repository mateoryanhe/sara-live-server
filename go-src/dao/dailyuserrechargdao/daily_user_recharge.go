package dailyuserrechargdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var dailyUserRechargeCacheMgr *cache.CacheMgr

func InitDailyUserRechargeDao() {
	dailyUserRechargeCacheMgr = cache.NewCacheMgr()
}

func getById(id string, date string, userId uint64) *entity.DailyUserRecharge {
	v := dailyUserRechargeCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.DailyUserRecharge
		_ = g.Model(string(entity.TbDailyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserRecharge(date, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.DailyUserRecharge)
	return row
}

// TryRecordRecharge 记录用户当日首次充值;本日已充值过返回 false
func TryRecordRecharge(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserRechargeId(date, userId)
	data := getById(id, date, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
