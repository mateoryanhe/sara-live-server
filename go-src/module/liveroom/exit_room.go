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
	//需要判断前端上传的房间号
	reqId := httpserver.GetReqId(ctx)
	user := userinfodao.GetUserInfoByUserId(userId)

	if reqId > 0 && user.LiveRoomId == req.RoomId && user.LiveRoomVer > reqId {
		return &liveroomdto.LeaveRoomRes{}, nil
	}
	exitRoom(userId, req.RoomId)
	return &liveroomdto.LeaveRoomRes{
		OnlineCount: getLenForRoom(req.RoomId),
	}, nil
}

func exitRoom(userId uint64, roomId uint64) {
	onlineId := entity.BuildLiveRoomOnlineId(userId, roomId)
	existing := liveroomdao.GetOnlineById(onlineId, userId, roomId)
	//wasAudienceOnline := userId != roomId && existing != nil &&
	//	existing.Status == entity.LiveRoomOnlineStatusOnline

	//先广播
	broadcastAudienceLeave(roomId, userId, getLenForRoom(roomId))
	if existing != nil && existing.Status != entity.LiveRoomOnlineStatusOffline {
		existing.SetStatus(entity.LiveRoomOnlineStatusOffline)
	}

	//if user := userinfodao.GetUserInfoByUserId(userId); user != nil && user.LiveRoomId == roomId {
	//	user.SetLiveRoomId(0)
	//	user.SetLiveRoomVer(0)
	//}
	removeOnline(userId, roomId)
	flushCommonOnlineMap(roomId)
}
