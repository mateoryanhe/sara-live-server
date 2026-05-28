package liveroom

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// JoinRoom 玩家加入直播间,记录状态置为 Online
func JoinRoom(ctx context.Context, req *liveroomdto.JoinRoomReq) (*liveroomdto.JoinRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)

	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	onlineId := entity.BuildLiveRoomOnlineId(userId, req.RoomId)
	existing := liveroomdao.GetOnlineById(onlineId)
	if existing == nil {
		o := entity.NewLiveRoomOnline(userId, req.RoomId)
		liveroomdao.AddOnlineToRoomCache(o)
	} else if existing.Status != entity.LiveRoomOnlineStatusOnline {
		existing.SetStatus(entity.LiveRoomOnlineStatusOnline)
		liveroomdao.AddOnlineToRoomCache(existing)
	}

	// 直播中且非主播本人:观众人数去重 +1
	if room.LiveRecordId > 0 && userId != req.RoomId {
		if liveroomdao.TryRecordLiveRecordAudience(room.LiveRecordId, userId) {
			if liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId); liveRecord != nil {
				liveRecord.AddTotalAudience(1)
			}
		}
	}

	return &liveroomdto.JoinRoomRes{
		OnlineId:    onlineId,
		OnlineCount: len(liveroomdao.GetOnlinesByRoom(req.RoomId)),
	}, nil
}
