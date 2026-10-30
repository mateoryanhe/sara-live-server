package statdao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/entity/stat"
)

var dailyApkPromotionChannelStatCacheMgr *cache.RowCache[*entity.DailyApkPromotionChannelStat]

func initDailyApkPromotionChannelStatDao() {
	dailyApkPromotionChannelStatCacheMgr = cache.NewRowCache[*entity.DailyApkPromotionChannelStat]()
}

func GetDailyApkPromotionChannelStat(date, channelCode string) *entity.DailyApkPromotionChannelStat {
	id := entity.BuildDailyApkPromotionChannelStatId(date, channelCode)
	return dailyApkPromotionChannelStatCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.DailyApkPromotionChannelStat, error) {
		var data *entity.DailyApkPromotionChannelStat
		_ = g.Model(string(entity.TbDailyApkPromotionChannelStat)).Unscoped().Where(g.Map{
			string(db.IdName): id,
		}).Scan(&data)
		if data != nil {
			syncDailyApkPromotionChannelStatName(data)
			return data, nil
		}
		name := apkPromotionChannelName(channelCode)
		return entity.NewDailyApkPromotionChannelStat(date, channelCode, name), nil
	})
}

func syncDailyApkPromotionChannelStatName(stat *entity.DailyApkPromotionChannelStat) {
	if stat == nil || strings.TrimSpace(stat.ChannelCode) == "" {
		return
	}
	name := apkPromotionChannelName(stat.ChannelCode)
	if name != "" && stat.ChannelName != name {
		stat.SetChannelName(name)
	}
}

func apkPromotionChannelName(channelCode string) string {
	row := apkpromotionchanneldao.GetByChannelCode(channelCode)
	if row == nil {
		return ""
	}
	return row.ChannelName
}

func PublishDailyApkPromotionChannelStat(data *entity.DailyApkPromotionChannelStat) {
	if data == nil || data.ID == "" || dailyApkPromotionChannelStatCacheMgr == nil {
		return
	}
	dailyApkPromotionChannelStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
