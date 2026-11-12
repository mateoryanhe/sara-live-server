package call

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity/call"
	"xr-game-server/errercode"
)

// RepushIncomingCall 当前用户若有「呼叫中」来电订单,重新推送 LiveRoomCallRequest,便于恢复通话流程.
// 无来电时 calling=false(不报错); 响应同时带回推送字段,WS 未连上时也可拉起 UI.
func RepushIncomingCall(ctx context.Context, _ *calldto.RepushIncomingCallReq) (*calldto.RepushIncomingCallRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	order := findActiveCallingOrderForReceiver(userId)
	if order == nil {
		return &calldto.RepushIncomingCallRes{Calling: false}, nil
	}
	channelName := buildCallChannelName(order.CallerId, order.ReceiverId)
	item, err := pushCallRequest(order, channelName)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return &calldto.RepushIncomingCallRes{Calling: false}, nil
	}
	return &calldto.RepushIncomingCallRes{
		Calling:        true,
		OrderId:        item.OrderId,
		Source:         item.Source,
		CallerId:       item.CallerId,
		ChannelName:    item.ChannelName,
		CallType:       item.CallType,
		CallerNickname: item.CallerNickname,
		CallerAvatar:   item.CallerAvatar,
		ReceiverToken:  item.ReceiverToken,
		Message:        item.Message,
	}, nil
}

// findActiveCallingOrderForReceiver 在活跃通话监控集中找「呼叫中」且接听方=userId 的订单.
func findActiveCallingOrderForReceiver(userId uint64) *entity.CallOrder {
	if userId == 0 {
		return nil
	}
	var matched *entity.CallOrder
	for _, orderId := range activeCallOrderIds.Slice() {
		order := calldao.GetOrderFromCache(orderId)
		if order == nil {
			order = calldao.GetOrderById(orderId)
		}
		if order == nil || order.ReceiverId != userId || !order.IsCalling() {
			continue
		}
		if matched == nil || order.CallStartTime.After(matched.CallStartTime) {
			matched = order
		}
	}
	return matched
}
