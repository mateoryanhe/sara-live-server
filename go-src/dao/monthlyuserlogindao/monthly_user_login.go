package monthlyuserlogindao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var monthlyUserLoginCacheMgr *cache.CacheMgr

func InitMonthlyUserLoginDao() {
	monthlyUserLoginCacheMgr = cache.NewCacheMgr()
}

func getById(id string, month string, userId uint64) *entity.MonthlyUserLogin {
	v := monthlyUserLoginCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.MonthlyUserLogin
		_ = g.Model(string(entity.TbMonthlyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserLogin(month, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.MonthlyUserLogin)
	return row
}

// TryRecordLogin 记录用户当月首次登录;已登录过返回 false
func TryRecordLogin(month string, userId uint64) bool {
	id := entity.BuildMonthlyUserLoginId(month, userId)
	data := getById(id, month, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	return true
}
