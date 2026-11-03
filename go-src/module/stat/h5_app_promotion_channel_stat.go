package stat

import (
	"time"

	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/module/promotionattribution"
)

func recordDailyH5AppPromotionChannelRegister(at time.Time, channelCode int) {
	if channelCode <= 0 || promotionchannelurldao.GetByChannelCode(channelCode) == nil {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyH5AppPromotionChannelStat(date, channelCode)
	stat.AddRegisterCount(1)
	statdao.PublishDailyH5AppPromotionChannelStat(stat)
}

func recordDailyH5AppPromotionChannelOrder(at time.Time, userId uint64, price float64) {
	if !userEligibleForPromotionAppDailyStat(userId) {
		return
	}
	channelCode := userinfodao.GetPromotionChannelCode(userId)
	apkCh := userinfodao.GetApkPromotionChannel(userId)
	if channelCode <= 0 || price <= 0 || !promotionattribution.UserHasValidPair(channelCode, apkCh) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyH5AppPromotionChannelStat(date, channelCode)
	stat.AddOrderCount(1)
	stat.AddOrderTotalAmount(price)
	statdao.PublishDailyH5AppPromotionChannelStat(stat)
}

func recordDailyH5AppPromotionChannelRechargeSuccess(at time.Time, userId uint64, amount float64) {
	if userId == 0 || amount <= 0 {
		return
	}
	if !userEligibleForPromotionAppDailyStat(userId) {
		return
	}
	channelCode := userinfodao.GetPromotionChannelCode(userId)
	apkCh := userinfodao.GetApkPromotionChannel(userId)
	if channelCode <= 0 || !promotionattribution.UserHasValidPair(channelCode, apkCh) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyH5AppPromotionChannelStat(date, channelCode)
	stat.AddRechargeSuccessAmount(amount)
	if statdao.TryRecordDailyH5AppPromotionChannelPayUser(date, channelCode, userId) {
		stat.AddPayUserCount(1)
	}
	statdao.PublishDailyH5AppPromotionChannelStat(stat)
}
