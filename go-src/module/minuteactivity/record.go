package minuteactivity

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/statdao"
	"xr-game-server/gameevent"
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
	metricCall1v1ConnectSuccess
	metricCall1v1Consume
	metricRechargeOrder
)

func schedule(kind metricKind, at time.Time, platform string) {
	if at.IsZero() {
		at = time.Now()
	}
	platform = gameevent.NormalizeStatClientPlatform(platform)
	kindCopy := kind
	atCopy := at
	platformCopy := platform
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		record(kindCopy, atCopy, platformCopy)
	})
}

func record(kind metricKind, at time.Time, platform string) {
	platform = gameevent.NormalizeStatClientPlatform(platform)
	row := statdao.GetMinuteLiveActivityStatByTimePlatform(at, platform)
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
	case metricCall1v1ConnectSuccess:
		row.AddCall1v1ConnectSuccessCount(1)
	case metricCall1v1Consume:
		row.AddCall1v1ConsumeCount(1)
	case metricRechargeOrder:
		row.AddRechargeOrderCount(1)
	default:
		return
	}
	statdao.PublishMinuteLiveActivityStat(at, platform, row)
}

// RecordGift 直播间送礼成功
func RecordGift(at time.Time) {
	schedule(metricGift, at, "")
}

// RecordGiftFromCtx 带 X-Platform 的送礼统计
func RecordGiftFromCtx(ctx context.Context, at time.Time) {
	schedule(metricGift, at, httpserver.ClientPlatformFromContext(ctx))
}

// RecordChat 直播间公屏文字或付费弹幕
func RecordChat(at time.Time) {
	schedule(metricChat, at, "")
}

// RecordChatFromCtx 带 X-Platform 的聊天统计
func RecordChatFromCtx(ctx context.Context, at time.Time) {
	schedule(metricChat, at, httpserver.ClientPlatformFromContext(ctx))
}

func onDayCleanup(_ any) {
	cutoff := statentity.TruncateUTCMinute(time.Now().AddDate(0, 0, -retentionDays))
	statdao.DeleteMinuteLiveActivityBefore(cutoff)
}
