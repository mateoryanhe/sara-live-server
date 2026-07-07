package call

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
)

// AcceptCall 接收者同意接听通话
func AcceptCall(ctx context.Context, req *calldto.AcceptCallReq) (*calldto.AcceptCallRes, error) {
	receiverId := httpserver.GetAuthId(ctx)
	if receiverId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if err := ensureNotInCall(receiverId); err != nil {
		return nil, err
	}

	order := calldao.GetOrderById(req.OrderId)
	if order == nil {
		return nil, errercode.CreateCode(errercode.CallOrderNonExist)
	}

	if order.CallStatus != entity.CallOrderStatusCalling {
		return nil, errercode.CreateCode(errercode.CallOrderStateInvalid)
	}

	channelName := buildCallChannelName(order.CallerId, order.ReceiverId)
	token, tokenExpireAt, err := resolveChannelToken(receiverId, channelName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if err := chargeLiveRoomCallOnAccept(order, now); err != nil {
		return nil, err
	}
	order.SetCallStatus(entity.CallOrderStatusInCall)
	order.SetAnswerTime(&now)
	calldao.FlushOrderCache(order)

	upsertCallUser(receiverId, order.ID, now)
	if callerUser := calldao.GetUserById(order.CallerId); callerUser != nil {
		callerUser.SetHeartTime(&now)
		calldao.FlushUserCache(callerUser)
	}

	pushCallAccepted(order.CallerId, receiverId, order.ID, channelName, order.CallType)

	appId := ""
	if agoraCfg, err := agora.GetAppId(ctx, nil); err == nil && agoraCfg != nil {
		appId = agoraCfg.AppId
	}

	return &calldto.AcceptCallRes{
		Success:     true,
		OrderId:     strconv.FormatUint(order.ID, 10),
		ChannelName: channelName,
		Token:       token,
		AppId:       appId,
		UserAccount: strconv.FormatUint(receiverId, 10),
		ExpireAt:    tokenExpireAt.Unix(),
	}, nil
}

func upsertCallUser(userId, orderId uint64, now time.Time) {
	callUser := calldao.GetUserById(userId)
	if callUser == nil {
		callUser = entity.NewCallUser(userId, orderId)
	} else {
		callUser.SetCallOrderId(orderId)
	}
	callUser.SetHeartTime(&now)
	calldao.FlushUserCache(callUser)
}
