package messagedao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const preloadUnreadDetailLimit = 100

// PreloadMessageToCache 批量预热消息相关 DAO 缓存
func PreloadMessageToCache(userIds []uint64) {
	GetAllActivityMessagesCached()
	if len(userIds) == 0 {
		return
	}
	for _, userId := range userIds {
		if userId == 0 {
			continue
		}
		preloadUserMessageCaches(userId)
	}
}

func preloadUserMessageCaches(userId uint64) {
	GetUnReadByUserId(userId)
	GetSystemMessageUnreadListCache(userId)
	GetUserActivityMessageListCacheA(userId)
	GetUserPersonalSystemMessageListCache(userId)
	preloadUnreadDetailsByUserId(userId)
}

func preloadUnreadDetailsByUserId(userId uint64) {
	list := make([]*entity.UserMessageUnreadDetail, 0)
	_ = g.Model(string(entity.TbUserMessageUnreadDetail)).
		Where("user_id = ?", userId).
		Order("updated_at desc").
		Limit(preloadUnreadDetailLimit).
		Scan(&list)
	for _, row := range list {
		FlushUnreadDetailCache(row)
	}
}
