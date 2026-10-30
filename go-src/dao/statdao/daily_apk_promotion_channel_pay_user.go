package statdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var dailyApkPromotionChannelPayUserCacheMgr *cache.RowCache[*entity.DailyApkPromotionChannelPayUser]

func initDailyApkPromotionChannelPayUserDao() {
	dailyApkPromotionChannelPayUserCacheMgr = cache.NewRowCache[*entity.DailyApkPromotionChannelPayUser]()
}

func getDailyApkPromotionChannelPayUserById(id, date, channelCode string, userId uint64) *entity.DailyApkPromotionChannelPayUser {
	return dailyApkPromotionChannelPayUserCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyApkPromotionChannelPayUser, error) {
		var row *entity.DailyApkPromotionChannelPayUser
		_ = g.Model(string(entity.TbDailyApkPromotionChannelPayUser)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewDailyApkPromotionChannelPayUser(date, channelCode, userId), nil
		}
		return row, nil
	})
}

func TryRecordDailyApkPromotionChannelPayUser(date, channelCode string, userId uint64) bool {
	if userId == 0 || date == "" || channelCode == "" {
		return false
	}
	id := entity.BuildDailyApkPromotionChannelPayUserId(date, channelCode, userId)
	data := getDailyApkPromotionChannelPayUserById(id, date, channelCode, userId)
	if data == nil || data.CreatedAt != nil {
		return false
	}
	now := time.Now()
	data.SetCreatedAt(&now)
	publishDailyApkPromotionChannelPayUser(data)
	return true
}

func publishDailyApkPromotionChannelPayUser(data *entity.DailyApkPromotionChannelPayUser) {
	if data == nil || data.ID == "" || dailyApkPromotionChannelPayUserCacheMgr == nil {
		return
	}
	dailyApkPromotionChannelPayUserCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
