package call

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const callAnswerConfirmChargeThreshold uint32 = 2

// ConfirmCall 通话应答确认;双方各确认一次,均确认后发起首次扣费
func ConfirmCall(ctx context.Context, req *calldto.ConfirmCallReq) (*calldto.ConfirmCallRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	lockKey := fmt.Sprintf("call_confirm_%d", req.OrderId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}
	if order.CallerId != userId && order.ReceiverId != userId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !order.IsAccepted() {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	if order.AnswerConfirmCount() >= callAnswerConfirmChargeThreshold {
		return buildConfirmCallRes(order, false), nil
	}
	if order.HasUserConfirmed(userId) {
		return buildConfirmCallRes(order, false), nil
	}

	firstChargeExecuted, err := applyCallAnswerConfirm(order, userId, time.Now())
	if err != nil {
		return nil, err
	}
	calldao.FlushOrderCache(order)

	return buildConfirmCallRes(order, firstChargeExecuted), nil
}

func applyCallAnswerConfirm(order *entity.CallOrder, userId uint64, now time.Time) (bool, error) {
	otherConfirmed := (userId == order.CallerId && order.ReceiverConfirmTime != nil) ||
		(userId == order.ReceiverId && order.CallerConfirmTime != nil)

	order.SetUserConfirmTime(userId, now)
	if !otherConfirmed {
		return false, nil
	}
	if order.ChargeTime != nil {
		return false, nil
	}

	if err := checkLiveRoomCallDiamondOnAccept(order); err != nil {
		clearUserConfirmTime(order, userId)
		return false, err
	}
	if err := chargeLiveRoomCallOnAccept(order, now); err != nil {
		clearUserConfirmTime(order, userId)
		return false, err
	}
	order.SetStatus(entity.CallOrderStatusInCall)
	pushCallStarted(order, now.Unix())
	return true, nil
}

func clearUserConfirmTime(order *entity.CallOrder, userId uint64) {
	if order == nil {
		return
	}
	if userId == order.CallerId {
		order.SetCallerConfirmTime(nil)
		return
	}
	if userId == order.ReceiverId {
		order.SetReceiverConfirmTime(nil)
	}
}

func buildConfirmCallRes(order *entity.CallOrder, firstChargeExecuted bool) *calldto.ConfirmCallRes {
	return &calldto.ConfirmCallRes{
		Success:             true,
		OrderId:             strconv.FormatUint(order.ID, 10),
		AnswerConfirmCount:  order.AnswerConfirmCount(),
		FirstChargeExecuted: firstChargeExecuted,
	}
}
