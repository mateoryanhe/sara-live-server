package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserLoginCacheMgr *cache.RowCache[*entity.DailyUserLogin]

func initDailyUserLoginDao() {
	dailyUserLoginCacheMgr = cache.NewRowCache[*entity.DailyUserLogin]()
}

func getDailyUserLoginById(id string, date string, userId uint64) *entity.DailyUserLogin {
	v := dailyUserLoginCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserLogin, error) {
		var row *entity.DailyUserLogin
		_ = g.Model(string(entity.TbDailyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserLogin(date, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordDailyLogin 记录用户当日首次登录;已登录过返回 false
func TryRecordDailyLogin(date string, userId uint64) bool {
	id := entity.BuildDailyUserLoginId(date, userId)
	data := getDailyUserLoginById(id, date, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishDailyUserLogin(data)
	return true
}
