package promotionchannelurl

import (
	"context"
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dto/promotionchannelurldto"
	"xr-game-server/errercode"
)

const maxPromotionDailyStatRangeDays = 92

// GetDailyStatList CMS 推广渠道日统计(按日期范围)
func GetDailyStatList(_ context.Context, req *promotionchannelurldto.PromotionChannelDailyStatListReq) (*promotionchannelurldto.PromotionChannelDailyStatListRes, error) {
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
	maxSpan := time.Duration(maxPromotionDailyStatRangeDays-1) * 24 * time.Hour
	if end.Sub(start) > maxSpan {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := statdao.ListDailyPromotionChannelStats(req)
	return &promotionchannelurldto.PromotionChannelDailyStatListRes{
		Total: total,
		Data:  list,
	}, nil
}

// GetH5AppPromotionDailyStatList CMS H5下载App推广日统计
func GetH5AppPromotionDailyStatList(_ context.Context, req *promotionchannelurldto.H5AppPromotionDailyStatListReq) (*promotionchannelurldto.H5AppPromotionDailyStatListRes, error) {
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
	maxSpan := time.Duration(maxPromotionDailyStatRangeDays-1) * 24 * time.Hour
	if end.Sub(start) > maxSpan {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	total, list := statdao.ListDailyH5AppPromotionChannelStats(req)
	return &promotionchannelurldto.H5AppPromotionDailyStatListRes{
		Total: total,
		Data:  list,
	}, nil
}
