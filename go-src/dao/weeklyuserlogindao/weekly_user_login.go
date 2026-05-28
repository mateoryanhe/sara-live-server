package weeklyuserlogindao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var weeklyUserLoginCacheMgr *cache.CacheMgr

func InitWeeklyUserLoginDao() {
	weeklyUserLoginCacheMgr = cache.NewCacheMgr()
}

func getById(id string, week string, userId uint64) *entity.WeeklyUserLogin {
	v := weeklyUserLoginCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.WeeklyUserLogin
		_ = g.Model(string(entity.TbWeeklyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserLogin(week, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.WeeklyUserLogin)
	return row
}

// TryRecordLogin 记录用户当周首次登录;已登录过返回 false
func TryRecordLogin(week string, userId uint64) bool {
	id := entity.BuildWeeklyUserLoginId(week, userId)
	data := getById(id, week, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
