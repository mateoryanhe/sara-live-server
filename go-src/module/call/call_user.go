package call

import (
	"time"

	"xr-game-server/dao/calldao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const callActiveHeartInterval = 30 * time.Second

// finishCallOrderIfHeartTimeout 心跳超过间隔未更新时,结束关联通话订单
func finishCallOrderIfHeartTimeout(callUser *entity.CallUser) {
	if callUser == nil || callUser.HeartTime == nil || callUser.CallOrderId == 0 {
		return
	}
	if time.Since(*callUser.HeartTime) <= callActiveHeartInterval {
		return
	}

	order := calldao.GetOrderById(callUser.CallOrderId)
	if order == nil || order.HasEnded() {
		return
	}

	now := time.Now()
	if order.IsCallStarted() {
		_ = refundLiveRoomCallLastMinuteIfNeeded(order, now)
	}
	callUser.SetHeartTime(nil)
	callUser.SetCallOrderId(0)
	if order.ReceiverId == callUser.ID {
		order.SetReceiverHangUpTime(&now)
	}
	if order.CallerId == callUser.ID {
		order.SetCallerHangUpTime(&now)
	}
	if !order.HasEnded() {
		order.SetOrderEndTime(&now)
		order.SetStatus(entity.CallOrderStatusEnded)
	}
	calldao.FlushOrderCache(order)
}

// ensureNotInCall 校验用户是否正在通话中
func ensureNotInCall(userId uint64) error {
	callUser := calldao.GetUserById(userId)
	if callUser == nil || callUser.CallOrderId == 0 || callUser.HeartTime == nil {
		return nil
	}

	finishCallOrderIfHeartTimeout(callUser)

	if time.Since(*callUser.HeartTime) < callActiveHeartInterval {
		return errercode.CreateCode(errercode.CallUserInCall)
	}
	return nil
}
