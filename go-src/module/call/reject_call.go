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

// AnchorRejectCall 主播拒接通话,支持不是主播,也可以通话
func AnchorRejectCall(ctx context.Context, req *calldto.AnchorRejectCallReq) (*calldto.AnchorRejectCallRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	if anchorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}
	if order.CallStatus != entity.CallOrderStatusCalling {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	order.SetCallStatus(entity.CallOrderStatusRejected)
	order.SetReceiverHangUpTime(&now)
	order.SetOrderEndTime(&now)
	calldao.FlushOrderCache(order)

	clearCallAnswerConfirmState(order.ID)
	resetCallUser(order.CallerId)
	pushCallRejected(order.CallerId, anchorId, order.ID)

	return &calldto.AnchorRejectCallRes{Success: true}, nil
}

func resetCallUser(userId uint64) {
	callUser := calldao.GetUserById(userId)
	if callUser == nil {
		return
	}
	callUser.SetHeartTime(nil)
	callUser.SetCallOrderId(0)
	calldao.FlushUserCache(callUser)
}
