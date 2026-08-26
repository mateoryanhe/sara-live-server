package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyUserRechargeCacheMgr *cache.RowCache[*entity.DailyUserRecharge]

func initDailyUserRechargeDao() {
	dailyUserRechargeCacheMgr = cache.NewRowCache[*entity.DailyUserRecharge]()
}

func getDailyUserRechargeById(id string, date string, userId uint64) *entity.DailyUserRecharge {
	v := dailyUserRechargeCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyUserRecharge, error) {
		var row *entity.DailyUserRecharge
		_ = g.Model(string(entity.TbDailyUserRecharge)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyUserRecharge(date, userId), nil
		}
		return row, nil
	})
	return v
}

// TryRecordDailyRecharge 记录用户当日首次充值;本日已充值过返回 false
func TryRecordDailyRecharge(date string, userId uint64) bool {
	if userId == 0 || date == "" {
		return false
	}
	id := entity.BuildDailyUserRechargeId(date, userId)
	data := getDailyUserRechargeById(id, date, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishDailyUserRecharge(data)
	return true
}
