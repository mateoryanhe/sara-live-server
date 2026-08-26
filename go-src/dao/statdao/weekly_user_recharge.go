package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var weeklyUserRechargeCacheMgr *cache.RowCache[*entity.WeeklyUserRecharge]

func initWeeklyUserRechargeDao() {
	weeklyUserRechargeCacheMgr = cache.NewRowCache[*entity.WeeklyUserRecharge]()
}

func getWeeklyUserRechargeById(id string, week string, userId uint64) *entity.WeeklyUserRecharge {
	v := weeklyUserRechargeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.WeeklyUserRecharge, error) {
		var row *entity.WeeklyUserRecharge
		_ = g.Model(string(entity.TbWeeklyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewWeeklyUserRecharge(week, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordWeeklyRecharge 记录用户当周首次充值;本周已充值过返回 false
func TryRecordWeeklyRecharge(week string, userId uint64) bool {
	if userId == 0 || week == "" {
		return false
	}
	id := entity.BuildWeeklyUserRechargeId(week, userId)
	data := getWeeklyUserRechargeById(id, week, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishWeeklyUserRecharge(data)
	return true
}
