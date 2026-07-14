package call

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/calldao"
	"xr-game-server/entity"
)

const callRingTimeout = 60 * time.Second

var activeCallOrderIds = gset.NewTSet[uint64](true)

func initCallOrderWatch() {
	restoreActiveCallOrders()
	xrtimer.AddSingleton(gctx.New(), time.Second, callOrderWatchTick)
}

func restoreActiveCallOrders() {
	orders := calldao.LoadUnclosedOrders()
	if len(orders) == 0 {
		return
	}

	userIds := gset.NewTSet[uint64](true)
	for _, order := range orders {
		if order == nil || order.ID == 0 || order.HasEnded() {
			continue
		}
		calldao.AddOrderToCache(order)
		trackActiveCallOrder(order.ID)
		userIds.Add(order.CallerId)
		userIds.Add(order.ReceiverId)
	}
	calldao.PreloadCallUsersToCache(userIds.Slice())

	ctx := gctx.New()
	g.Log().Infof(ctx, "恢复通话订单监控,count=%v", activeCallOrderIds.Size())
}

func trackActiveCallOrder(orderId uint64) {
	if orderId == 0 {
		return
	}
	activeCallOrderIds.Add(orderId)
}

func untrackActiveCallOrder(orderId uint64) {
	if orderId == 0 {
		return
	}
	activeCallOrderIds.Remove(orderId)
}

func callOrderWatchTick(_ context.Context) {
	now := time.Now()
	for _, orderId := range activeCallOrderIds.Slice() {
		order := calldao.GetOrderFromCache(orderId)
		if order == nil || order.HasEnded() {
			untrackActiveCallOrder(orderId)
			continue
		}
		if processAbnormalCallOrder(order, now) {
			untrackActiveCallOrder(orderId)
		}
	}
}

func processAbnormalCallOrder(order *entity.CallOrder, now time.Time) bool {
	if order == nil || order.HasEnded() {
		return false
	}
	if order.IsCalling() {
		if now.Sub(order.CallStartTime) >= callRingTimeout {
			finishCallOrderOnRingTimeout(order, now)
			return true
		}
		caller := calldao.GetUserFromCache(order.CallerId)
		if isCallUserHeartOffline(caller, order.ID, now) {
			finishCallOrderOnHeartTimeout(order, now)
			return true
		}
		return false
	}
	if order.HasAnswered() {
		if shouldFinishCallOrderOnHeartTimeout(order, now) {
			finishCallOrderOnHeartTimeout(order, now)
			return true
		}
		if order.IsCallStarted() {
			if err := chargeLiveRoomCallBillingIfDue(order, now); err != nil && isDiamondNotEnough(err) {
				finishCallOrderOnBillingFailed(order, now)
				return true
			}
		}
	}
	return false
}

func shouldFinishCallOrderOnHeartTimeout(order *entity.CallOrder, now time.Time) bool {
	caller := calldao.GetUserFromCache(order.CallerId)
	receiver := calldao.GetUserFromCache(order.ReceiverId)
	return isCallUserHeartOffline(caller, order.ID, now) ||
		isCallUserHeartOffline(receiver, order.ID, now)
}

func isCallUserHeartOffline(callUser *entity.CallUser, orderId uint64, now time.Time) bool {
	if callUser == nil || callUser.CallOrderId != orderId || callUser.HeartTime == nil || callUser.HeartTime.IsZero() {
		return true
	}
	return now.Sub(*callUser.HeartTime) > callActiveHeartInterval
}

func finishCallOrderOnRingTimeout(order *entity.CallOrder, now time.Time) {
	if order == nil || order.HasEnded() {
		return
	}
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusCallTimeout)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)
	pushCallTimeout(order.CallerId, order.ReceiverId, order.ID)
	untrackActiveCallOrder(order.ID)
}

func finishCallOrderOnHeartTimeout(order *entity.CallOrder, now time.Time) {
	if order == nil || order.HasEnded() {
		return
	}
	if order.IsCallStarted() {
		_ = refundLiveRoomCallLastMinuteIfNeeded(order, now)
	}

	endUserId := resolveCallOfflineUserId(order, now)
	if endUserId == order.CallerId {
		order.SetCallerHangUpTime(&now)
	} else if endUserId == order.ReceiverId {
		order.SetReceiverHangUpTime(&now)
	}
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusHeartTimeout)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)
	pushCallEnded(order.CallerId, endUserId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)
	pushCallEnded(order.ReceiverId, endUserId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)
	untrackActiveCallOrder(order.ID)
}

func resolveCallOfflineUserId(order *entity.CallOrder, now time.Time) uint64 {
	caller := calldao.GetUserFromCache(order.CallerId)
	receiver := calldao.GetUserFromCache(order.ReceiverId)
	if isCallUserHeartOffline(caller, order.ID, now) {
		return order.CallerId
	}
	if isCallUserHeartOffline(receiver, order.ID, now) {
		return order.ReceiverId
	}
	return order.CallerId
}
