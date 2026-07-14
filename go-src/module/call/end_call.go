package call

import (
	"context"
	"errors"
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
	if order.HasEnded() {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}
	if !order.IsCalling() && !order.HasAnswered() {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	if order.IsCallStarted() {
		if err := refundLiveRoomCallLastMinuteIfNeeded(order, now); err != nil {
			return nil, err
		}
	}
	if userId == order.CallerId {
		order.SetCallerHangUpTime(&now)
	} else {
		order.SetReceiverHangUpTime(&now)
	}
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusEnded)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)

	peerId := order.ReceiverId
	if userId == order.ReceiverId {
		peerId = order.CallerId
	}
	pushCallEnded(peerId, userId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)
	untrackActiveCallOrder(order.ID)

	return &calldto.EndCallRes{
		Success:         true,
		OrderId:         strconv.FormatUint(order.ID, 10),
		CallDuration:    order.CallDuration,
		BillingDuration: order.BillingDuration,
		TotalCost:       order.TotalCost,
	}, nil
}

// finishCallOrderOnBillingFailed 钻石不足续费时结束通话,并通知双方
func finishCallOrderOnBillingFailed(order *entity.CallOrder, now time.Time) {
	if order == nil || order.HasEnded() {
		return
	}
	if order.IsCallStarted() {
		_ = refundLiveRoomCallLastMinuteIfNeeded(order, now)
	}
	order.SetCallerHangUpTime(&now)
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusBillingFailed)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)

	endUserId := order.CallerId
	pushCallEndedDueToBillingFailed(order.ReceiverId, endUserId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)
	pushCallEndedDueToBillingFailed(order.CallerId, endUserId, order.ID, order.CallDuration, order.BillingDuration, order.TotalCost)
	untrackActiveCallOrder(order.ID)
}

func isDiamondNotEnough(err error) bool {
	var bizErr *errercode.XError
	return errors.As(err, &bizErr) && bizErr.Code() == errercode.DiamondNotEnough
}
