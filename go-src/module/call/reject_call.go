package call

import (
	"context"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity/call"
	"xr-game-server/errercode"
)

// AnchorRejectCall 目标用户拒接通话；保留既有函数名兼容接口。
func AnchorRejectCall(ctx context.Context, req *calldto.AnchorRejectCallReq) (*calldto.AnchorRejectCallRes, error) {
	receiverId := httpserver.GetAuthId(ctx)
	if receiverId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}
	if order.ReceiverId != receiverId {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if !order.IsCalling() {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	now := time.Now()
	order.SetReceiverHangUpTime(&now)
	order.SetOrderEndTime(&now)
	order.SetStatus(entity.CallOrderStatusRejected)
	calldao.FlushOrderCache(order)

	resetCallUser(order.CallerId)
	pushCallRejected(order)
	untrackActiveCallOrder(order.ID)

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
