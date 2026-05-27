package dailyuserlogindao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var dailyUserLoginCacheMgr *cache.CacheMgr

func InitDailyUserLoginDao() {
	dailyUserLoginCacheMgr = cache.NewCacheMgr()
}

func GetById(id string, date string, userId uint64) *entity.DailyUserLogin {
	v := dailyUserLoginCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.DailyUserLogin
		_ = g.Model(string(entity.TbDailyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserLogin(date, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.DailyUserLogin)
	return row
}

// TryRecordLogin 记录用户当日首次登录;已登录过返回 false
func TryRecordLogin(date string, userId uint64) bool {
	id := entity.BuildDailyUserLoginId(date, userId)
	data := GetById(id, date, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
