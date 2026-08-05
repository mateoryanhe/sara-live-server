package liveroom

import (
	"xr-game-server/core/push"
)

// PushToRoomAudience 向直播间全部在线观众推送.
func PushToRoomAudience(roomId uint64, pushCmd int, payload any) {
	if roomId == 0 {
		return
	}
	for _, userId := range getOnline(roomId) {
		push.Data(userId, pushCmd, payload)
	}
}
