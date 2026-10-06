package call

import (
	"time"

	"xr-game-server/dao/liveroomdao"
	callentity "xr-game-server/entity/call"
	"xr-game-server/module/livecfg"
)

func canSkipOneToOneDiamondPrecheck(payerId, anchorId uint64) bool {
	if payerId == 0 || anchorId == 0 {
		return false
	}
	if livecfg.GetOneToOneDailyFreeSeconds() == 0 {
		return false
	}
	return liveroomdao.IsOneToOneDailyFreeAvailable(payerId, anchorId)
}

func canSkipOneToOneDiamondPrecheckForOrder(order *callentity.CallOrder) bool {
	if order == nil || order.Source != callentity.CallOrderSourceOneToOneRoom {
		return false
	}
	anchorId, payerId, billable := callDiamondBillingPartyIds(order)
	if !billable {
		return false
	}
	return canSkipOneToOneDiamondPrecheck(payerId, anchorId)
}

// tryApplyOneToOneDailyFreeOnAccept 接通后延迟首分钟扣费并作废当日整段免费额度。
func tryApplyOneToOneDailyFreeOnAccept(order *callentity.CallOrder, anchorId, payerId uint64, now time.Time) bool {
	if order == nil || order.Source != callentity.CallOrderSourceOneToOneRoom {
		return false
	}
	freeSec := livecfg.GetOneToOneDailyFreeSeconds()
	if freeSec == 0 || payerId == 0 || anchorId == 0 {
		return false
	}
	if !liveroomdao.IsOneToOneDailyFreeAvailable(payerId, anchorId) {
		return false
	}
	liveroomdao.MarkOneToOneDailyFreeUsed(payerId, anchorId)
	order.SetOneToOneFreeDeferSeconds(freeSec)
	if order.PricePerMinute > 0 {
		nextCharge := now.Add(time.Duration(freeSec) * time.Second)
		order.SetChargeTime(&nextCharge)
	}
	return true
}
