package stat

import (
	"time"

	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/module/promotionattribution"
)

func recordDailyApkPromotionChannelRegister(at time.Time, channelCode string) {
	code, ok := promotionattribution.NormalizeApkPromotionChannelCode(channelCode)
	if !ok {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyApkPromotionChannelStat(date, code)
	stat.AddRegisterCount(1)
	statdao.PublishDailyApkPromotionChannelStat(stat)
}

func recordDailyApkPromotionChannelDownloadClick(at time.Time, channelCode string) {
	code, ok := promotionattribution.NormalizeApkPromotionChannelCode(channelCode)
	if !ok {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyApkPromotionChannelStat(date, code)
	stat.AddDownloadClickCount(1)
	statdao.PublishDailyApkPromotionChannelStat(stat)
}

func recordDailyApkPromotionChannelOrder(at time.Time, userId uint64, price float64) {
	if userId == 0 || price <= 0 {
		return
	}
	channelCode, ok := promotionattribution.NormalizeApkPromotionChannelCode(userinfodao.GetApkPromotionChannel(userId))
	if !ok {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyApkPromotionChannelStat(date, channelCode)
	stat.AddOrderCount(1)
	stat.AddOrderTotalAmount(price)
	statdao.PublishDailyApkPromotionChannelStat(stat)
}

func recordDailyApkPromotionChannelRechargeSuccess(at time.Time, userId uint64, amount float64) {
	if userId == 0 || amount <= 0 {
		return
	}
	channelCode, ok := promotionattribution.NormalizeApkPromotionChannelCode(userinfodao.GetApkPromotionChannel(userId))
	if !ok {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	stat := statdao.GetDailyApkPromotionChannelStat(date, channelCode)
	stat.AddRechargeSuccessAmount(amount)
	if statdao.TryRecordDailyApkPromotionChannelPayUser(date, channelCode, userId) {
		stat.AddPayUserCount(1)
	}
	statdao.PublishDailyApkPromotionChannelStat(stat)
}

// EnqueueApkPromotionChannelDownloadClick 渠道短链 GET 302 前入队统计.
func EnqueueApkPromotionChannelDownloadClick(channelCode string, at time.Time) {
	code, ok := promotionattribution.NormalizeApkPromotionChannelCode(channelCode)
	if !ok {
		return
	}
	enqueue(&statJob{Kind: jobApkPromotionChannelDownload, Payload: &apkPromotionChannelDownloadPayload{
		ChannelCode: code,
		StatAt:      at,
	}})
}

type apkPromotionChannelDownloadPayload struct {
	ChannelCode string
	StatAt      time.Time
}

func consumeApkPromotionChannelDownloadJob(job *statJob) {
	payload, ok := job.Payload.(*apkPromotionChannelDownloadPayload)
	if !ok || payload == nil || payload.ChannelCode == "" {
		return
	}
	at := payload.StatAt
	if at.IsZero() {
		at = time.Now()
	}
	recordDailyApkPromotionChannelDownloadClick(at, payload.ChannelCode)
}
