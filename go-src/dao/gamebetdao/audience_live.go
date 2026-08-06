package gamebetdao

import (
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
)

// ResolveAudienceLiveContext 获取观众当前所在直播间与直播记录 ID.
func ResolveAudienceLiveContext(userId uint64) (liveRoomId, liveRecordId uint64) {
	if userId == 0 {
		return 0, 0
	}
	userInfo := userinfodao.GetUserInfoByUserId(userId)
	if userInfo == nil || userInfo.LiveRoomId == 0 {
		return 0, 0
	}
	liveRoomId = userInfo.LiveRoomId
	room := liveroomdao.GetRoomById(liveRoomId)
	if room == nil || room.LiveRecordId == 0 {
		return liveRoomId, 0
	}
	return liveRoomId, room.LiveRecordId
}
