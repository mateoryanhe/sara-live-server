package liveroomdao

import (
	"time"

	"xr-game-server/entity/live"
)

func withDailyGuildEffectiveLive(guildId uint64, at time.Time, fn func(*entity.DailyGuildEffectiveLive)) {
	if guildId == 0 || fn == nil {
		return
	}
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyGuildEffectiveLive(date, guildId)
	if row == nil {
		return
	}
	fn(row)
	PublishDailyGuildEffectiveLive(row)
}

// MirrorDailyGuildGiftEarn 同步礼物收益到工会日表
func MirrorDailyGuildGiftEarn(guildId uint64, at time.Time, amount float64) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.AddGiftEarn(amount)
	})
}

// MirrorDailyGuildPaidDanmakuEarn 同步付费弹幕收益到工会日表
func MirrorDailyGuildPaidDanmakuEarn(guildId uint64, at time.Time, amount float64) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.AddPaidDanmakuEarn(amount)
	})
}

// MirrorDailyGuildPrivateRoomTicketEarn 同步私密房门票收益到工会日表
func MirrorDailyGuildPrivateRoomTicketEarn(guildId uint64, at time.Time, amount float64) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.AddPrivateRoomTicketEarn(amount)
	})
}

// MirrorDailyGuildPrivateRoomWatchEarn 同步私密房观看收益到工会日表
func MirrorDailyGuildPrivateRoomWatchEarn(guildId uint64, at time.Time, amount float64) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.AddPrivateRoomWatchEarn(amount)
	})
}

// MirrorDailyGuildVideoCallIncomeDelta 同步通话收益增减到工会日表
func MirrorDailyGuildVideoCallIncomeDelta(guildId uint64, at time.Time, amount float64, ticket, billing bool) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.ApplyVideoCallIncomeDelta(amount, ticket, billing)
	})
}

// MirrorDailyGuildLiveDuration 同步心跳上报直播时长到工会日表
func MirrorDailyGuildLiveDuration(guildId uint64, at time.Time, sec float64) {
	withDailyGuildEffectiveLive(guildId, at, func(row *entity.DailyGuildEffectiveLive) {
		row.AddTotalLiveDuration(sec)
	})
}
