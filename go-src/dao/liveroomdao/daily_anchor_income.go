package liveroomdao

import (
	"time"

	"xr-game-server/entity/live"
)

func withDailyAnchorEffectiveLive(roomId uint64, at time.Time, fn func(*entity.DailyAnchorEffectiveLive)) {
	if roomId == 0 || fn == nil {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyAnchorEffectiveLive(date, roomId)
	if row == nil {
		return
	}
	fn(row)
}

// MirrorDailyAnchorGiftEarn 同步礼物收益到主播日表
func MirrorDailyAnchorGiftEarn(roomId uint64, at time.Time, amount float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddGiftEarn(amount)
	})
}

// MirrorDailyAnchorPaidDanmakuEarn 同步付费弹幕收益到主播日表
func MirrorDailyAnchorPaidDanmakuEarn(roomId uint64, at time.Time, amount float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddPaidDanmakuEarn(amount)
	})
}

// MirrorDailyAnchorPrivateRoomTicketEarn 同步私密房门票收益到主播日表
func MirrorDailyAnchorPrivateRoomTicketEarn(roomId uint64, at time.Time, amount float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddPrivateRoomTicketEarn(amount)
	})
}

// MirrorDailyAnchorPrivateRoomWatchEarn 同步私密房观看收益到主播日表
func MirrorDailyAnchorPrivateRoomWatchEarn(roomId uint64, at time.Time, amount float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddPrivateRoomWatchEarn(amount)
	})
}

// MirrorDailyAnchorShortVideoEarn 同步短视频付费观看收益到主播日表
func MirrorDailyAnchorShortVideoEarn(roomId uint64, at time.Time, amount float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddShortVideoEarn(amount)
	})
}

// MirrorDailyAnchorGameEarn 同步游戏收益到主播日表
func MirrorDailyAnchorGameEarn(roomId uint64, at time.Time, goldAmount, incomeDelta float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddGameEarn(goldAmount, incomeDelta)
	})
}

// MirrorDailyAnchorVideoCallIncomeDelta 同步通话收益增减到主播日表
func MirrorDailyAnchorVideoCallIncomeDelta(roomId uint64, at time.Time, amount float64, ticket, billing bool) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.ApplyVideoCallIncomeDelta(amount, ticket, billing)
	})
}

// MirrorDailyAnchorLiveDuration 同步心跳上报直播时长到主播日表
func MirrorDailyAnchorLiveDuration(roomId uint64, at time.Time, sec float64) {
	withDailyAnchorEffectiveLive(roomId, at, func(row *entity.DailyAnchorEffectiveLive) {
		row.AddTotalLiveDuration(sec)
	})
}
