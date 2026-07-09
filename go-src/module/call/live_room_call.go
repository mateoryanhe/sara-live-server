package call

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/calldao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
	"xr-game-server/module/wallet"
)

func buildCallChannelName(callerId, receiverId uint64) string {
	return fmt.Sprintf("%d_%d", callerId, receiverId)
}

// LiveRoomCall 直播间向主播发起通话呼叫
func LiveRoomCall(ctx context.Context, req *calldto.LiveRoomCallReq) (*calldto.LiveRoomCallRes, error) {
	callerId := httpserver.GetAuthId(ctx)
	anchorId := req.AnchorId
	if callerId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if callerId == anchorId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := ensureNotInCall(callerId); err != nil {
		return nil, err
	}
	if err := ensureNotInCall(anchorId); err != nil {
		return nil, err
	}

	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}

	requiredDiamond := room.Ticket
	if room.Billing > 0 {
		requiredDiamond += room.Billing
	}
	if requiredDiamond > 0 {
		if err := wallet.DiamondNotEnough(callerId, requiredDiamond); err != nil {
			return nil, err
		}
	}

	order := entity.NewCallOrder(
		callerId,
		anchorId,
		entity.CallOrderTypeVideo,
		entity.CallOrderSourceLiveRoom,
		strconv.FormatUint(room.LiveRecordId, 10),
		room.Ticket,
		room.Billing,
	)
	calldao.AddOrderToCache(order)

	channelName := buildCallChannelName(callerId, anchorId)
	if err := pushLiveRoomCallRequest(anchorId, callerId, order.ID, channelName, order.CallType); err != nil {
		return nil, err
	}

	token, tokenExpireAt, err := agora.ResolveCallToken(callerId, channelName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	callUser := entity.NewCallUser(callerId, order.ID)
	callUser.SetHeartTime(&now)
	calldao.AddUserToCache(callUser)

	agoraCfg, _ := agora.GetAppId(ctx, nil)
	appId := ""
	if agoraCfg != nil {
		appId = agoraCfg.AppId
	}

	return &calldto.LiveRoomCallRes{
		OrderId:     strconv.FormatUint(order.ID, 10),
		ChannelName: channelName,
		Token:       token,
		AppId:       appId,
		UserAccount: strconv.FormatUint(callerId, 10),
		ExpireAt:    tokenExpireAt,
	}, nil
}
