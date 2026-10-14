package statdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/entity/stat"
)

var dailyPromotionChannelStatCacheMgr *cache.RowCache[*entity.DailyPromotionChannelStat]

func initDailyPromotionChannelStatDao() {
	dailyPromotionChannelStatCacheMgr = cache.NewRowCache[*entity.DailyPromotionChannelStat]()
}

// GetDailyPromotionChannelStat 按日期与渠道编码获取日统计,不存在则新建
func GetDailyPromotionChannelStat(date string, channelCode int) *entity.DailyPromotionChannelStat {
	id := entity.BuildDailyPromotionChannelStatId(date, channelCode)
	return dailyPromotionChannelStatCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyPromotionChannelStat, error) {
		var data *entity.DailyPromotionChannelStat
		_ = g.Model(string(entity.TbDailyPromotionChannelStat)).Unscoped().Where(g.Map{
			string(db.IdName): id,
		}).Scan(&data)
		if data != nil {
			syncDailyPromotionChannelStatName(data)
			return data, nil
		}
		name := promotionChannelName(channelCode)
		return entity.NewDailyPromotionChannelStat(date, channelCode, name), nil
	})
}

func syncDailyPromotionChannelStatName(stat *entity.DailyPromotionChannelStat) {
	if stat == nil || stat.ChannelCode <= 0 {
		return
	}
	name := promotionChannelName(stat.ChannelCode)
	if name != "" && stat.ChannelName != name {
		stat.SetChannelName(name)
	}
}

func promotionChannelName(channelCode int) string {
	if channelCode <= 0 {
		return ""
	}
	row := promotionchannelurldao.GetByChannelCode(channelCode)
	if row == nil {
		return ""
	}
	return row.ChannelName
}

func PublishDailyPromotionChannelStat(data *entity.DailyPromotionChannelStat) {
	if data == nil || data.ID == "" || dailyPromotionChannelStatCacheMgr == nil {
		return
	}
	dailyPromotionChannelStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
