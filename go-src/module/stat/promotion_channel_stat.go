package stat

import (
	"time"

	"xr-game-server/dao/promotionchannelurldao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	rechargeentity "xr-game-server/entity/recharge"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func recordDailyPromotionChannelRegister(at time.Time, channelCode int) {
	if channelCode <= 0 || promotionchannelurldao.GetByChannelCode(channelCode) == nil {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyPromotionChannelStat(date, channelCode)
	stat.AddRegisterCount(1)
	statdao.PublishDailyPromotionChannelStat(stat)
}

func recordDailyPromotionChannelH5DownloadApk(at time.Time, channelCode int) {
	if channelCode <= 0 || promotionchannelurldao.GetByChannelCode(channelCode) == nil {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyPromotionChannelStat(date, channelCode)
	stat.AddH5DownloadApkCount(1)
	statdao.PublishDailyPromotionChannelStat(stat)
}

func recordDailyPromotionChannelOrder(at time.Time, userId uint64, price float64) {
	channelCode := userinfodao.GetPromotionChannelCode(userId)
	if channelCode <= 0 || price <= 0 {
		return
	}
	if promotionchannelurldao.GetByChannelCode(channelCode) == nil {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyPromotionChannelStat(date, channelCode)
	stat.AddOrderCount(1)
	stat.AddOrderTotalAmount(price)
	statdao.PublishDailyPromotionChannelStat(stat)
}

func recordDailyPromotionChannelRechargeSuccess(at time.Time, userId uint64, amount float64) {
	if userId == 0 || amount <= 0 {
		return
	}
	channelCode := userinfodao.GetPromotionChannelCode(userId)
	if channelCode <= 0 {
		return
	}
	if promotionchannelurldao.GetByChannelCode(channelCode) == nil {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyPromotionChannelStat(date, channelCode)
	stat.AddRechargeSuccessAmount(amount)
	if statdao.TryRecordDailyPromotionChannelPayUser(date, channelCode, userId) {
		stat.AddPayUserCount(1)
	}
	statdao.PublishDailyPromotionChannelStat(stat)
}

func consumePromotionChannelOrderJob(job *statJob) {
	val, ok := job.Payload.(*gameevent.RechargeOrderCreatedEventData)
	if !ok || val == nil || val.OrderId == 0 {
		return
	}
	order := rechargeorderdao.GetById(val.OrderId)
	if order == nil || order.UserId == 0 || order.Price <= 0 {
		return
	}
	if order.Source == rechargeentity.RechargeOrderSourceManual ||
		order.Source == rechargeentity.RechargeOrderSourceCoinMerchant {
		return
	}
	at := order.CreatedAt
	if at.IsZero() {
		at = time.Now()
	}
	recordDailyPromotionChannelOrder(at, order.UserId, order.Price)
	recordDailyApkPromotionChannelOrder(at, order.UserId, order.Price)
}
