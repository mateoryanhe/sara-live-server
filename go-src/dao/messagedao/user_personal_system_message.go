package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
)

const UserPersonalSystemMessageListCacheMax = 50

var userPersonalSystemMessageListCacheMgr *cache.CacheMgr

func initUserPersonalSystemMessageDao() {
	userPersonalSystemMessageListCacheMgr = cache.NewCacheMgr()
}

// GetUserPersonalSystemMessageListCache 获取用户个人系统消息列表缓存(前50条,按创建时间倒序)
func GetUserPersonalSystemMessageListCache(userId uint64) []*entity.UserPersonalSystemMessage {
	if userId == 0 || userPersonalSystemMessageListCacheMgr == nil {
		return nil
	}
	v := userPersonalSystemMessageListCacheMgr.GetData(userId, func(ctx context.Context) (interface{}, error) {
		list := make([]*entity.UserPersonalSystemMessage, 0)
		_ = g.Model(string(entity.TbUserPersonalSystemMessage)).Ctx(context.Background()).
			Where("user_id = ?", userId).
			Order("created_at desc").
			Limit(UserPersonalSystemMessageListCacheMax).
			Scan(&list)
		return list, nil
	})
	list, _ := v.([]*entity.UserPersonalSystemMessage)
	return list
}

// FlushUserPersonalSystemMessageListCache 刷新用户个人系统消息列表缓存
func FlushUserPersonalSystemMessageListCache(userId uint64, list []*entity.UserPersonalSystemMessage) {
	if userId == 0 || userPersonalSystemMessageListCacheMgr == nil {
		return
	}
	if list == nil {
		list = make([]*entity.UserPersonalSystemMessage, 0)
	}
	userPersonalSystemMessageListCacheMgr.FlushCache(userId, list)
}
