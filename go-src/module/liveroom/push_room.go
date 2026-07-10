package liveroom

import (
	"xr-game-server/core/push"
)

// PushToRoomAudience 向直播间在线观众推送(不含 excludeUserIds 中的用户)
func PushToRoomAudience(roomId uint64, pushCmd int, payload any, excludeUserIds ...uint64) {
	if roomId == 0 {
		return
	}
	exclude := make(map[uint64]struct{}, len(excludeUserIds))
	for _, userId := range excludeUserIds {
		if userId == 0 {
			continue
		}
		exclude[userId] = struct{}{}
	}
	for _, userId := range getOnline(roomId) {
		if _, skip := exclude[userId]; skip {
			continue
		}
		push.Data(userId, pushCmd, payload)
	}
}
