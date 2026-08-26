package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserLoginCacheMgr *cache.RowCache[*entity.WeeklyUserLogin]

func initWeeklyUserLoginDao() {
	weeklyUserLoginCacheMgr = cache.NewRowCache[*entity.WeeklyUserLogin]()
}

func getWeeklyUserLoginById(id string, week string, userId uint64) *entity.WeeklyUserLogin {
	v := weeklyUserLoginCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserLogin, error) {
		var row *entity.WeeklyUserLogin
		_ = g.Model(string(entity.TbWeeklyUserLogin)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserLogin(week, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordWeeklyLogin 记录用户当周首次登录;已登录过返回 false
func TryRecordWeeklyLogin(week string, userId uint64) bool {
	id := entity.BuildWeeklyUserLoginId(week, userId)
	data := getWeeklyUserLoginById(id, week, userId)
	if data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishWeeklyUserLogin(data)
	return true
}
