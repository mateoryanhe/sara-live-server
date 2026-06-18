package liveroom

import (
	"context"

	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/stat"
)

// JoinRoom 玩家加入直播间,记录状态置为 Online
func JoinRoom(ctx context.Context, req *liveroomdto.JoinRoomReq) (*liveroomdto.JoinRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)

	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	if userId != room.ID && room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}

	now := time.Now()
	onlineId := entity.BuildLiveRoomOnlineId(userId, req.RoomId)
	existing := liveroomdao.GetOnlineById(onlineId, userId, room.ID)
	if existing.IsKickBanned() {
		return nil, errercode.CreateCode(errercode.LiveRoomKickBanned)
	}

	if err := ensureCanJoinPrivateRoom(userId, room); err != nil {
		return nil, err
	}

	ticketDeducted, err := chargePrivateRoomTicketIfNeeded(userId, room, now)
	if err != nil {
		return nil, err
	}

	//alreadyOnline := existing.Status == entity.LiveRoomOnlineStatusOnline
	existing.SetStatus(entity.LiveRoomOnlineStatusOnline)
	existing.SetJoinTime(&now)
	existing.SetHeartTime(&now)
	addToOnline(userId, room.ID)

	if userId != room.ID {
		if user := userinfodao.GetUserInfoByUserId(userId); user != nil {
			user.SetLiveRoomId(room.ID)
		}
		broadcastAudienceJoin(room.ID, userId, getLenForRoom(room.ID))
	}

	// 直播中且非主播本人:有效观众(单场去重 + 日/周/月跨直播间去重)
	if room.LiveRecordId > 0 && userId != req.RoomId {
		stat.RecordValidAudience(userId, time.Now())
		if liveroomdao.TryRecordLiveRecordAudience(room.LiveRecordId, userId) {
			if liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId); liveRecord != nil {
				liveRecord.AddTotalAudience(1)
			}
		}
	}

	return &liveroomdto.JoinRoomRes{
		OnlineId:       onlineId,
		OnlineCount:    getLenForRoom(room.ID),
		TicketDeducted: ticketDeducted,
	}, nil
}
