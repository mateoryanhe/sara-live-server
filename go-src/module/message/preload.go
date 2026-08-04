package message

import "xr-game-server/dao/messagedao"

// PreloadMessageToCache 批量预热消息模块缓存(含 DAO 层与私信未读列表)
func PreloadMessageToCache(userIds []uint64) {
	messagedao.PreloadMessageToCache(userIds)
	if len(userIds) == 0 {
		return
	}
	for _, userId := range userIds {
		if userId == 0 {
			continue
		}
		getPrivateMessageUnreadList(userId)
	}
}
