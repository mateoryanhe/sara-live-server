package call

import (
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

const liveRoomCallRequestMessage = "有人请求跟你通话"

func pushLiveRoomCallRequest(receiverId, callerId, orderId uint64, channelName string, callType uint8) error {
	receiverToken, _, err := resolveChannelToken(receiverId, channelName)
	if err != nil {
		return err
	}

	item := &calldto.CallRequestPushItem{
		OrderId:       strconv.FormatUint(orderId, 10),
		CallerId:      strconv.FormatUint(callerId, 10),
		ChannelName:   channelName,
		CallType:      callType,
		ReceiverToken: receiverToken,
		Message:       liveRoomCallRequestMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(callerId); u != nil {
		item.CallerNickname = u.Nickname
		item.CallerAvatar = upload.ResolveAvatarUrlForUser(callerId, u.Avatar)
	}
	push.Data(receiverId, cmd.LiveRoomCallRequest, item)
	return nil
}

const liveRoomCallRejectedMessage = "对方已拒接"

func pushCallRejected(callerId, receiverId, orderId uint64) {
	item := &calldto.CallRejectedPushItem{
		OrderId:    strconv.FormatUint(orderId, 10),
		ReceiverId: strconv.FormatUint(receiverId, 10),
		Message:    liveRoomCallRejectedMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(receiverId); u != nil {
		item.ReceiverNickname = u.Nickname
		item.ReceiverAvatar = upload.ResolveAvatarUrlForUser(receiverId, u.Avatar)
	}
	push.Data(callerId, cmd.LiveRoomCallRejected, item)
}

const liveRoomCallAcceptedMessage = "对方已接听"

func pushCallAccepted(callerId, receiverId, orderId uint64, channelName string, callType uint8) {
	item := &calldto.CallAcceptedPushItem{
		OrderId:     strconv.FormatUint(orderId, 10),
		ReceiverId:  strconv.FormatUint(receiverId, 10),
		ChannelName: channelName,
		CallType:    callType,
		Message:     liveRoomCallAcceptedMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(receiverId); u != nil {
		item.ReceiverNickname = u.Nickname
		item.ReceiverAvatar = upload.ResolveAvatarUrlForUser(receiverId, u.Avatar)
	}
	push.Data(callerId, cmd.LiveRoomCallAccepted, item)
}

const liveRoomCallEndedMessage = "通话已结束"

func pushCallEnded(peerId, endUserId, orderId uint64, callDuration, billingDuration uint32, totalCost float64) {
	item := &calldto.CallEndedPushItem{
		OrderId:         strconv.FormatUint(orderId, 10),
		EndUserId:       strconv.FormatUint(endUserId, 10),
		CallDuration:    callDuration,
		BillingDuration: billingDuration,
		TotalCost:       totalCost,
		Message:         liveRoomCallEndedMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(endUserId); u != nil {
		item.EndUserNickname = u.Nickname
		item.EndUserAvatar = upload.ResolveAvatarUrlForUser(endUserId, u.Avatar)
	}
	push.Data(peerId, cmd.LiveRoomCallEnded, item)
}

const liveRoomCallStartedMessage = "通话已开始"

func pushCallStarted(order *entity.CallOrder, startedAt int64) {
	if order == nil {
		return
	}
	channelName := buildCallChannelName(order.CallerId, order.ReceiverId)
	item := &calldto.CallStartedPushItem{
		OrderId:     strconv.FormatUint(order.ID, 10),
		CallerId:    strconv.FormatUint(order.CallerId, 10),
		ReceiverId:  strconv.FormatUint(order.ReceiverId, 10),
		ChannelName: channelName,
		CallType:    order.CallType,
		StartedAt:   startedAt,
		Message:     liveRoomCallStartedMessage,
	}
	push.Data(order.CallerId, cmd.LiveRoomCallStarted, item)
	push.Data(order.ReceiverId, cmd.LiveRoomCallStarted, item)
}

const liveRoomCallTimeoutMessage = "呼叫超时"

func pushCallTimeout(callerId, receiverId, orderId uint64) {
	item := &calldto.CallTimeoutPushItem{
		OrderId:    strconv.FormatUint(orderId, 10),
		CallerId:   strconv.FormatUint(callerId, 10),
		ReceiverId: strconv.FormatUint(receiverId, 10),
		Message:    liveRoomCallTimeoutMessage,
	}
	push.Data(callerId, cmd.LiveRoomCallTimeout, item)
	push.Data(receiverId, cmd.LiveRoomCallTimeout, item)
}
