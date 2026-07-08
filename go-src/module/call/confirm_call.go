package call

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const callAnswerConfirmChargeThreshold uint32 = 2

var callAnswerConfirmFirstUser sync.Map // orderId -> userId

// ConfirmCall 通话应答确认;双方各确认一次,累计达到2次时发起首次扣费
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
	if order.CallStatus != entity.CallOrderStatusInCall {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	if order.AnswerConfirmCount >= callAnswerConfirmChargeThreshold {
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
	switch order.AnswerConfirmCount {
	case 0:
		callAnswerConfirmFirstUser.Store(order.ID, userId)
		order.AddAnswerConfirmCount(1)
		return false, nil
	case 1:
		if firstUser, ok := callAnswerConfirmFirstUser.Load(order.ID); ok && firstUser.(uint64) == userId {
			return false, errercode.CreateCode(errercode.CallOrderStateInvalid)
		}
		if err := checkLiveRoomCallDiamondOnAccept(order); err != nil {
			return false, err
		}
		if err := chargeLiveRoomCallOnAccept(order, now); err != nil {
			return false, err
		}
		order.AddAnswerConfirmCount(1)
		clearCallAnswerConfirmState(order.ID)
		pushCallStarted(order, now.Unix())
		return true, nil
	default:
		return false, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}
}

func buildConfirmCallRes(order *entity.CallOrder, firstChargeExecuted bool) *calldto.ConfirmCallRes {
	return &calldto.ConfirmCallRes{
		Success:             true,
		OrderId:             strconv.FormatUint(order.ID, 10),
		AnswerConfirmCount:  order.AnswerConfirmCount,
		FirstChargeExecuted: firstChargeExecuted,
	}
}

func clearCallAnswerConfirmState(orderId uint64) {
	callAnswerConfirmFirstUser.Delete(orderId)
}
