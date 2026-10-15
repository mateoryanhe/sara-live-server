package call

import (
	"time"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	callentity "xr-game-server/entity/call"
	userentity "xr-game-server/entity/user"
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

func isOrdinaryViewer(userId uint64) bool {
	if userId == 0 {
		return false
	}
	info := userinfodao.GetUserInfoByUserId(userId)
	return info != nil && info.UserType == userentity.UserTypeNormal
}

func orderSupportsDailyFreeCall(order *callentity.CallOrder) bool {
	if order == nil {
		return false
	}
	switch order.Source {
	case callentity.CallOrderSourceOneToOneRoom:
		return true
	case callentity.CallOrderSourceLiveRoom:
		return isOrdinaryViewer(order.CallerId)
	default:
		return false
	}
}

func canSkipOneToOneDiamondPrecheckForOrder(order *callentity.CallOrder) bool {
	if !orderSupportsDailyFreeCall(order) {
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
	if !orderSupportsDailyFreeCall(order) {
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
