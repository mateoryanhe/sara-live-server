package liveroom

import "xr-game-server/entity"

// markLiveRoomCreated App端创建/完善直播间后,标记 user_infos.has_live_room = true
func markLiveRoomCreated(user *entity.UserInfo) {
	if user == nil || user.HasLiveRoom {
		return
	}
	user.SetHasLiveRoom(true)
}
