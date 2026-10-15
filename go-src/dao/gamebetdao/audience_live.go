package gamebetdao

import (
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
)

// ResolveAudienceLiveContext 游戏下注关联直播间(user_ext.game_live_room_id; 0=大厅,不写 live_room_id/live_record_id).
func ResolveAudienceLiveContext(userId uint64) (liveRoomId, liveRecordId uint64) {
	if userId == 0 {
		return 0, 0
	}
	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil || ext.GameLiveRoomId == 0 {
		return 0, 0
	}
	liveRoomId = ext.GameLiveRoomId
	room := liveroomdao.GetRoomById(liveRoomId)
	if room == nil || room.LiveRecordId == 0 {
		return liveRoomId, 0
	}
	return liveRoomId, room.LiveRecordId
}
