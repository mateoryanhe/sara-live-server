package call

import (
	"time"

	"xr-game-server/dao/calldao"
	"xr-game-server/entity/call"
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

	finishCallOrderOnHeartTimeout(order, time.Now())
}

// ensureNotInCall 校验用户是否正在通话中; busyCode 区分「本人忙线」与「对方忙线」等提示.
func ensureNotInCall(userId uint64, busyCode errercode.XRCode) error {
	if busyCode == 0 {
		busyCode = errercode.CallUserInCall
	}
	callUser := calldao.GetUserById(userId)
	if callUser == nil || callUser.CallOrderId == 0 || callUser.HeartTime == nil {
		return nil
	}

	finishCallOrderIfHeartTimeout(callUser)

	if time.Since(*callUser.HeartTime) < callActiveHeartInterval {
		return errercode.CreateCode(busyCode)
	}
	return nil
}
