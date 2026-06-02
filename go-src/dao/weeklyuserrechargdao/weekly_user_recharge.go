package weeklyuserrechargdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var weeklyUserRechargeCacheMgr *cache.CacheMgr

func InitWeeklyUserRechargeDao() {
	weeklyUserRechargeCacheMgr = cache.NewCacheMgr()
}

func getById(id string, week string, userId uint64) *entity.WeeklyUserRecharge {
	v := weeklyUserRechargeCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.WeeklyUserRecharge
		_ = g.Model(string(entity.TbWeeklyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserRecharge(week, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.WeeklyUserRecharge)
	return row
}

// TryRecordRecharge 记录用户当周首次充值;本周已充值过返回 false
func TryRecordRecharge(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserRechargeId(week, userId)
	data := getById(id, week, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
