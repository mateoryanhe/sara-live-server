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

var dailyH5AppPromotionChannelStatCacheMgr *cache.RowCache[*entity.DailyH5AppPromotionChannelStat]

func initDailyH5AppPromotionChannelStatDao() {
	dailyH5AppPromotionChannelStatCacheMgr = cache.NewRowCache[*entity.DailyH5AppPromotionChannelStat]()
}

func GetDailyH5AppPromotionChannelStat(date string, channelCode int) *entity.DailyH5AppPromotionChannelStat {
	id := entity.BuildDailyH5AppPromotionChannelStatId(date, channelCode)
	return dailyH5AppPromotionChannelStatCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyH5AppPromotionChannelStat, error) {
		var data *entity.DailyH5AppPromotionChannelStat
		_ = g.Model(string(entity.TbDailyH5AppPromotionChannelStat)).Unscoped().Where(g.Map{
			string(db.IdName): id,
		}).Scan(&data)
		if data != nil {
			syncDailyH5AppPromotionChannelStatName(data)
			return data, nil
		}
		name := h5AppPromotionChannelName(channelCode)
		return entity.NewDailyH5AppPromotionChannelStat(date, channelCode, name), nil
	})
}

func syncDailyH5AppPromotionChannelStatName(stat *entity.DailyH5AppPromotionChannelStat) {
	if stat == nil || stat.ChannelCode <= 0 {
		return
	}
	name := h5AppPromotionChannelName(stat.ChannelCode)
	if name != "" && stat.ChannelName != name {
		stat.SetChannelName(name)
	}
}

func h5AppPromotionChannelName(channelCode int) string {
	row := promotionchannelurldao.GetByChannelCode(channelCode)
	if row == nil {
		return ""
	}
	return row.ChannelName
}

func PublishDailyH5AppPromotionChannelStat(data *entity.DailyH5AppPromotionChannelStat) {
	if data == nil || data.ID == "" || dailyH5AppPromotionChannelStatCacheMgr == nil {
		return
	}
	dailyH5AppPromotionChannelStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
