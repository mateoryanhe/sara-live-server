package call

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// CallTimeout 呼叫超时(仅呼叫中状态可处理)
func CallTimeout(ctx context.Context, req *calldto.CallTimeoutReq) (*calldto.CallTimeoutRes, error) {
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
		return &calldto.CallTimeoutRes{Success: true}, nil
	}
	if !order.IsCalling() {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusCallTimeout)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	resetCallUser(order.ReceiverId)
	pushCallTimeout(order.CallerId, order.ReceiverId, order.ID)

	return &calldto.CallTimeoutRes{Success: true}, nil
}
