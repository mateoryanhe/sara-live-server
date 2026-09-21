package call

import (
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity/call"
	"xr-game-server/module/agora"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/upload"
)

const liveRoomCallRequestMessage = "有人请求跟你通话"

func pushCallRequest(order *entity.CallOrder, channelName string) error {
	receiverToken, _, err := agora.ResolveChannelToken(order.ReceiverId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return err
	}

	item := &calldto.CallRequestPushItem{
		OrderId:       strconv.FormatUint(order.ID, 10),
		Source:        entity.NormalizeCallOrderSource(order.Source),
		CallerId:      strconv.FormatUint(order.CallerId, 10),
		ChannelName:   channelName,
		CallType:      order.CallType,
		ReceiverToken: receiverToken,
		Message:       liveRoomCallRequestMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(order.CallerId); u != nil {
		item.CallerNickname = u.Nickname
		item.CallerAvatar = upload.ResolveAvatarUrlForUser(order.CallerId, u.Avatar)
	}
	push.Data(order.ReceiverId, cmd.LiveRoomCallRequest, item)
	return nil
}

const liveRoomCallRejectedMessage = "对方已拒接"

func pushCallRejected(order *entity.CallOrder) {
	if order == nil {
		return
	}
	item := &calldto.CallRejectedPushItem{
		OrderId:    strconv.FormatUint(order.ID, 10),
		Source:     entity.NormalizeCallOrderSource(order.Source),
		ReceiverId: strconv.FormatUint(order.ReceiverId, 10),
		Message:    liveRoomCallRejectedMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(order.ReceiverId); u != nil {
		item.ReceiverNickname = u.Nickname
		item.ReceiverAvatar = upload.ResolveAvatarUrlForUser(order.ReceiverId, u.Avatar)
	}
	push.Data(order.CallerId, cmd.LiveRoomCallRejected, item)
}

const liveRoomCallAcceptedMessage = "对方已接听"

func pushCallAccepted(order *entity.CallOrder, channelName string) {
	if order == nil {
		return
	}
	item := &calldto.CallAcceptedPushItem{
		OrderId:     strconv.FormatUint(order.ID, 10),
		Source:      entity.NormalizeCallOrderSource(order.Source),
		ReceiverId:  strconv.FormatUint(order.ReceiverId, 10),
		ChannelName: channelName,
		CallType:    order.CallType,
		Message:     liveRoomCallAcceptedMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(order.ReceiverId); u != nil {
		item.ReceiverNickname = u.Nickname
		item.ReceiverAvatar = upload.ResolveAvatarUrlForUser(order.ReceiverId, u.Avatar)
	}
	push.Data(order.CallerId, cmd.LiveRoomCallAccepted, item)
}

const liveRoomCallAnchorAcceptedAudienceMessage = "主播开始接听视频通话"

func buildLiveRoomCallAnchorAcceptedAudienceItem(order *entity.CallOrder) *calldto.CallAnchorAcceptedAudiencePushItem {
	if order == nil || order.Source != entity.CallOrderSourceLiveRoom {
		return nil
	}
	if order.CallType != entity.CallOrderTypeVideo {
		return nil
	}

	roomId := order.ReceiverId
	item := &calldto.CallAnchorAcceptedAudiencePushItem{
		RoomId:   strconv.FormatUint(roomId, 10),
		AnchorId: strconv.FormatUint(roomId, 10),
		OrderId:  strconv.FormatUint(order.ID, 10),
		Source:   entity.NormalizeCallOrderSource(order.Source),
		CallerId: strconv.FormatUint(order.CallerId, 10),
		CallType: order.CallType,
		Message:  liveRoomCallAnchorAcceptedAudienceMessage,
	}
	if u := userinfodao.GetUserInfoByUserId(roomId); u != nil {
		item.AnchorNickname = u.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrlForUser(roomId, u.Avatar)
	}
	return item
}

func pushLiveRoomCallAcceptedToAudienceUser(userId uint64, order *entity.CallOrder) {
	if userId == 0 {
		return
	}
	item := buildLiveRoomCallAnchorAcceptedAudienceItem(order)
	if item == nil {
		return
	}
	push.Data(userId, cmd.LiveRoomCallAnchorAcceptedAudience, item)
}

func pushLiveRoomCallAcceptedToAudience(order *entity.CallOrder) {
	item := buildLiveRoomCallAnchorAcceptedAudienceItem(order)
	if item == nil {
		return
	}
	liveroom.PushToRoomAudience(order.ReceiverId, cmd.LiveRoomCallAnchorAcceptedAudience, item)
}

const liveRoomCallEndedMessage = "通话已结束"
const liveRoomCallBillingFailedMessage = "钻石不足,通话已结束"

func pushCallEnded(peerId, endUserId uint64, order *entity.CallOrder) {
	pushCallEndedWithMessage(peerId, endUserId, order, liveRoomCallEndedMessage)
}

func pushCallEndedDueToBillingFailed(peerId, endUserId uint64, order *entity.CallOrder) {
	pushCallEndedWithMessage(peerId, endUserId, order, liveRoomCallBillingFailedMessage)
}

func pushCallEndedWithMessage(peerId, endUserId uint64, order *entity.CallOrder, message string) {
	if order == nil {
		return
	}
	item := &calldto.CallEndedPushItem{
		OrderId:         strconv.FormatUint(order.ID, 10),
		Source:          entity.NormalizeCallOrderSource(order.Source),
		EndUserId:       strconv.FormatUint(endUserId, 10),
		CallDuration:    order.CallDuration,
		BillingDuration: order.BillingDuration,
		TotalCost:       order.TotalCost,
		Message:         message,
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
		Source:      entity.NormalizeCallOrderSource(order.Source),
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

func pushCallTimeout(order *entity.CallOrder) {
	if order == nil {
		return
	}
	item := &calldto.CallTimeoutPushItem{
		OrderId:    strconv.FormatUint(order.ID, 10),
		Source:     entity.NormalizeCallOrderSource(order.Source),
		CallerId:   strconv.FormatUint(order.CallerId, 10),
		ReceiverId: strconv.FormatUint(order.ReceiverId, 10),
		Message:    liveRoomCallTimeoutMessage,
	}
	push.Data(order.CallerId, cmd.LiveRoomCallTimeout, item)
	push.Data(order.ReceiverId, cmd.LiveRoomCallTimeout, item)
}
