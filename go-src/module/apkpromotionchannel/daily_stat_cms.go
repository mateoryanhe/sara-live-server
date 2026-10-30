package apkpromotionchannel

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/apkpromotionchanneldto"
	"xr-game-server/errercode"
)

const maxApkPromotionDailyStatRangeDays = 92

func GetDailyStatList(_ context.Context, req *apkpromotionchanneldto.ApkPromotionChannelDailyStatListReq) (*apkpromotionchanneldto.ApkPromotionChannelDailyStatListRes, error) {
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	end, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
	if err != nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if end.Before(start) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	maxSpan := time.Duration(maxApkPromotionDailyStatRangeDays-1) * 24 * time.Hour
	if end.Sub(start) > maxSpan {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := statdao.ListDailyApkPromotionChannelStats(req)
	return &apkpromotionchanneldto.ApkPromotionChannelDailyStatListRes{
		Total: total,
		Data:  list,
	}, nil
}
