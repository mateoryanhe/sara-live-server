package call

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity/call"
)

func initLiveRoomCallJoinNotify() {
	event.Sub(event.LiveRoomAudienceJoined, onLiveRoomAudienceJoined)
}

func onLiveRoomAudienceJoined(val any) {
	payload, ok := val.(*liveroomdto.AudienceJoinedEvent)
	if !ok || payload == nil || payload.RoomId == 0 || payload.UserId == 0 {
		return
	}
	if payload.UserId == payload.RoomId {
		return
	}
	order := getActiveLiveRoomVideoCallOrder(payload.RoomId)
	if order == nil {
		return
	}
	pushLiveRoomCallAcceptedToAudienceUser(payload.UserId, order)
}

func getActiveLiveRoomVideoCallOrder(anchorId uint64) *entity.CallOrder {
	if anchorId == 0 {
		return nil
	}
	callUser := calldao.GetUserById(anchorId)
	if callUser == nil || callUser.CallOrderId == 0 {
		return nil
	}
	finishCallOrderIfHeartTimeout(callUser)

	order := calldao.GetOrderById(callUser.CallOrderId)
	if order == nil || order.HasEnded() {
		return nil
	}
	if order.Source != entity.CallOrderSourceLiveRoom {
		return nil
	}
	if order.CallType != entity.CallOrderTypeVideo {
		return nil
	}
	if order.ReceiverId != anchorId {
		return nil
	}
	if !order.HasAnswered() {
		return nil
	}
	return order
}
