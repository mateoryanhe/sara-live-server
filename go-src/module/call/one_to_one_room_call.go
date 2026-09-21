package call

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/dto/calldto"
	callentity "xr-game-server/entity/call"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
	"xr-game-server/module/wallet"
)

// OneToOneRoomCall 向1v1房间中的目标用户发起视频通话。
// targetId 只表示接听方；主播房间和付费用户根据双方用户类型解析。
func OneToOneRoomCall(ctx context.Context, req *calldto.OneToOneRoomCallReq) (*calldto.OneToOneRoomCallRes, error) {
	callerId := httpserver.GetAuthId(ctx)
	targetId := req.TargetId
	if callerId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if callerId == targetId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := ensureNotInCall(callerId); err != nil {
		return nil, err
	}
	if err := ensureNotInCall(targetId); err != nil {
		return nil, err
	}

	callCtx, err := resolveOneToOneRoomCallContext(callerId, targetId)
	if err != nil {
		return nil, err
	}
	if callCtx.pricePerMinute > 0 {
		if err := wallet.CanPayWithGoldExchange(callCtx.audienceId, callCtx.pricePerMinute); err != nil {
			return nil, err
		}
	}

	order := callentity.NewCallOrder(
		callerId,
		targetId,
		callCtx.audienceId,
		callentity.CallOrderTypeVideo,
		callentity.CallOrderSourceOneToOneRoom,
		callCtx.orderParams,
		0,
		callCtx.pricePerMinute,
	)
	calldao.AddOrderToCache(order)
	trackActiveCallOrder(order.ID)

	channelName := buildCallChannelName(callerId, targetId)
	if err := pushCallRequest(order, channelName); err != nil {
		return nil, err
	}

	token, tokenExpireAt, err := agora.ResolveChannelToken(callerId, channelName, agoradto.RTCRolePublisher)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	order.SetUserHeartTime(callerId, now)
	callUser := callentity.NewCallUser(callerId, order.ID)
	callUser.SetHeartTime(&now)
	calldao.AddUserToCache(callUser)
	calldao.FlushOrderCache(order)

	agoraCfg, _ := agora.GetAppId(ctx, nil)
	appId := ""
	if agoraCfg != nil {
		appId = agoraCfg.AppId
	}

	return &calldto.OneToOneRoomCallRes{
		OrderId:     strconv.FormatUint(order.ID, 10),
		ChannelName: channelName,
		Token:       token,
		AppId:       appId,
		UserAccount: strconv.FormatUint(callerId, 10),
		ExpireAt:    tokenExpireAt,
	}, nil
}
