package minuteactivity

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
)

type metricKind int

const (
	metricJoin metricKind = iota + 1
	metricLeave
	metricGift
	metricChat
	metricBet
	metricCall1v1
	metricRechargeOrder
)

func schedule(kind metricKind, at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		record(kind, at)
	})
}

func record(kind metricKind, at time.Time) {
	row := statdao.GetMinuteLiveActivityStatByTime(at)
	switch kind {
	case metricJoin:
		row.AddJoinCount(1)
	case metricLeave:
		row.AddLeaveCount(1)
	case metricGift:
		row.AddGiftCount(1)
	case metricChat:
		row.AddChatCount(1)
	case metricBet:
		row.AddBetCount(1)
	case metricCall1v1:
		row.AddCall1v1Count(1)
	case metricRechargeOrder:
		row.AddRechargeOrderCount(1)
	default:
		return
	}
	statdao.PublishMinuteLiveActivityStat(at, row)
}

// RecordGift 直播间送礼成功
func RecordGift(at time.Time) {
	schedule(metricGift, at)
}

// RecordChat 直播间公屏文字或付费弹幕
func RecordChat(at time.Time) {
	schedule(metricChat, at)
}

func onDayCleanup(_ any) {
	cutoff := statentity.TruncateUTCMinute(time.Now().AddDate(0, 0, -retentionDays))
	statdao.DeleteMinuteLiveActivityBefore(cutoff)
}
