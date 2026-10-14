package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyPromotionChannelPayUserCacheMgr *cache.RowCache[*entity.DailyPromotionChannelPayUser]

func initDailyPromotionChannelPayUserDao() {
	dailyPromotionChannelPayUserCacheMgr = cache.NewRowCache[*entity.DailyPromotionChannelPayUser]()
}

func getDailyPromotionChannelPayUserById(id, date string, channelCode int, userId uint64) *entity.DailyPromotionChannelPayUser {
	return dailyPromotionChannelPayUserCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyPromotionChannelPayUser, error) {
		var row *entity.DailyPromotionChannelPayUser
		_ = g.Model(string(entity.TbDailyPromotionChannelPayUser)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyPromotionChannelPayUser(date, channelCode, userId), nil
		}
		return row, nil
	})
}

// TryRecordDailyPromotionChannelPayUser 记录当日该渠道下用户首次支付;已记录返回 false
func TryRecordDailyPromotionChannelPayUser(date string, channelCode int, userId uint64) bool {
	if userId == 0 || date == "" || channelCode <= 0 {
		return false
	}
	id := entity.BuildDailyPromotionChannelPayUserId(date, channelCode, userId)
	data := getDailyPromotionChannelPayUserById(id, date, channelCode, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishDailyPromotionChannelPayUser(data)
	return true
}

func publishDailyPromotionChannelPayUser(data *entity.DailyPromotionChannelPayUser) {
	if data == nil || data.ID == "" || dailyPromotionChannelPayUserCacheMgr == nil {
		return
	}
	dailyPromotionChannelPayUserCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
