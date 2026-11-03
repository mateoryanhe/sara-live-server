package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyH5AppPromotionChannelPayUserCacheMgr *cache.RowCache[*entity.DailyH5AppPromotionChannelPayUser]

func initDailyH5AppPromotionChannelPayUserDao() {
	dailyH5AppPromotionChannelPayUserCacheMgr = cache.NewRowCache[*entity.DailyH5AppPromotionChannelPayUser]()
}

func getDailyH5AppPromotionChannelPayUserById(id, date string, channelCode int, userId uint64) *entity.DailyH5AppPromotionChannelPayUser {
	return dailyH5AppPromotionChannelPayUserCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyH5AppPromotionChannelPayUser, error) {
		var row *entity.DailyH5AppPromotionChannelPayUser
		_ = g.Model(string(entity.TbDailyH5AppPromotionChannelPayUser)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyH5AppPromotionChannelPayUser(date, channelCode, userId), nil
		}
		return row, nil
	})
}

func TryRecordDailyH5AppPromotionChannelPayUser(date string, channelCode int, userId uint64) bool {
	if userId == 0 || date == "" || channelCode <= 0 {
		return false
	}
	id := entity.BuildDailyH5AppPromotionChannelPayUserId(date, channelCode, userId)
	data := getDailyH5AppPromotionChannelPayUserById(id, date, channelCode, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishDailyH5AppPromotionChannelPayUser(data)
	return true
}

func publishDailyH5AppPromotionChannelPayUser(data *entity.DailyH5AppPromotionChannelPayUser) {
	if data == nil || data.ID == "" || dailyH5AppPromotionChannelPayUserCacheMgr == nil {
		return
	}
	dailyH5AppPromotionChannelPayUserCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
