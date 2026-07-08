package call

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// EndCall 结束通话(呼叫者或接听者主动挂断)
func EndCall(ctx context.Context, req *calldto.EndCallReq) (*calldto.EndCallRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}
	if order.CallerId != userId && order.ReceiverId != userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if isCallOrderFinished(order.CallStatus) {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}
	if order.CallStatus != entity.CallOrderStatusCalling && order.CallStatus != entity.CallOrderStatusInCall {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	if order.CallStatus == entity.CallOrderStatusInCall {
		if err := refundLiveRoomCallLastMinuteIfNeeded(order, now); err != nil {
			return nil, err
		}
	}
	if userId == order.CallerId {
		order.SetCallerHangUpTime(&now)
	} else {
		order.SetReceiverHangUpTime(&now)
	}
	order.SetCallStatus(entity.CallOrderStatusEnded)
	order.SetOrderEndTime(&now)
	calldao.FlushOrderCache(order)

	clearCallAnswerConfirmState(order.ID)
	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)

	peerId := order.ReceiverId
	if userId == order.ReceiverId {
		peerId = order.CallerId
	}
	pushCallEnded(peerId, userId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)

	return &calldto.EndCallRes{
		Success:         true,
		OrderId:         strconv.FormatUint(order.ID, 10),
		CallDuration:    order.CallDuration,
		BillingDuration: order.BillingDuration,
		TotalCost:       order.TotalCost,
	}, nil
}
