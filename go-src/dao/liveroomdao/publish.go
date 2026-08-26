package liveroomdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

func PublishLiveRoomOnline(data *entity.LiveRoomOnline) {
	if data == nil || data.ID == "" || onlineCacheMgr == nil {
		return
	}
	onlineCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishLiveRecord(data *entity.LiveRecord) {
	if data == nil || data.ID == 0 || liveRecordCacheMgr == nil {
		return
	}
	liveRecordCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishLiveRoomBillingPay(data *entity.LiveRoomBillingPay) {
	if data == nil || data.ID == "" || liveRoomBillingPayCacheMgr == nil {
		return
	}
	liveRoomBillingPayCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishRevenueLog(data *entity.LiveRevenueLog) {
	if data == nil || data.ID == 0 || revenueLogCacheMgr == nil {
		return
	}
	revenueLogCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishDailyAnchorEffectiveLive(data *entity.DailyAnchorEffectiveLive) {
	if data == nil || data.ID == "" || dailyAnchorEffectiveLiveCacheMgr == nil {
		return
	}
	dailyAnchorEffectiveLiveCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishDailyGuildEffectiveLive(data *entity.DailyGuildEffectiveLive) {
	if data == nil || data.ID == "" || dailyGuildEffectiveLiveCacheMgr == nil {
		return
	}
	dailyGuildEffectiveLiveCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
