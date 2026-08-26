package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var monthlyUserLoginCacheMgr *cache.RowCache[*entity.MonthlyUserLogin]

func initMonthlyUserLoginDao() {
	monthlyUserLoginCacheMgr = cache.NewRowCache[*entity.MonthlyUserLogin]()
}

func getMonthlyUserLoginById(id string, month string, userId uint64) *entity.MonthlyUserLogin {
	v := monthlyUserLoginCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.MonthlyUserLogin, error) {
		var row *entity.MonthlyUserLogin
		_ = g.Model(string(entity.TbMonthlyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewMonthlyUserLogin(month, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordMonthlyLogin 记录用户当月首次登录;已登录过返回 false
func TryRecordMonthlyLogin(month string, userId uint64) bool {
	id := entity.BuildMonthlyUserLoginId(month, userId)
	data := getMonthlyUserLoginById(id, month, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishMonthlyUserLogin(data)
	return true
}
