package liveroom

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
)

// LeaveRoom 玩家离开直播间,记录状态置为 Offline(保留行)
func LeaveRoom(ctx context.Context, req *liveroomdto.LeaveRoomReq) (*liveroomdto.LeaveRoomRes, error) {
	userId := httpserver.GetAuthId(ctx)
	user := userinfodao.GetUserInfoByUserId(userId)
	exitRoom(userId, user.LiveRoomId)
	return &liveroomdto.LeaveRoomRes{
		OnlineCount: getLenForRoom(user.LiveRoomId),
	}, nil
}

func exitRoom(userId uint64, roomId uint64) {
	onlineId := entity.BuildLiveRoomOnlineId(userId, roomId)
	existing := liveroomdao.GetOnlineById(onlineId, userId, roomId)
	wasAudienceOnline := userId != roomId && existing != nil &&
		existing.Status == entity.LiveRoomOnlineStatusOnline

	removeOnline(userId, roomId)
	if existing != nil && existing.Status != entity.LiveRoomOnlineStatusOffline {
		existing.SetStatus(entity.LiveRoomOnlineStatusOffline)
	}

	if wasAudienceOnline {
		if user := userinfodao.GetUserInfoByUserId(userId); user != nil && user.LiveRoomId == roomId {
			user.SetLiveRoomId(0)
		}
		broadcastAudienceLeave(roomId, userId, getLenForRoom(roomId))
	}
}
