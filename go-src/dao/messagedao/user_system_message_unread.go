package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
)

const SystemMessageUnreadListCacheMax = 100

var systemMessageUnreadListCacheMgr *cache.CacheMgr

func initSystemMessageUnreadDao() {
	systemMessageUnreadListCacheMgr = cache.NewCacheMgr()
}

// GetSystemMessageUnreadListCache 获取用户系统消息未读列表(优先读缓存,未命中查库)
func GetSystemMessageUnreadListCache(userId uint64) []*entity.UserSystemMessageUnread {
	if userId == 0 || systemMessageUnreadListCacheMgr == nil {
		return nil
	}
	v := systemMessageUnreadListCacheMgr.GetData(userId, func(ctx context.Context) (interface{}, error) {
		list := make([]*entity.UserSystemMessageUnread, 0)
		_ = g.Model(string(entity.TbUserSystemMessageUnread)).Ctx(context.Background()).
			Where("user_id = ? AND unread_count > 0", userId).
			Order("updated_at desc").
			Limit(SystemMessageUnreadListCacheMax).
			Scan(&list)
		return list, nil
	})
	list, _ := v.([]*entity.UserSystemMessageUnread)
	return list
}

// FlushSystemMessageUnreadListCache 刷新用户系统消息未读列表缓存
func FlushSystemMessageUnreadListCache(userId uint64, list []*entity.UserSystemMessageUnread) {
	if userId == 0 || systemMessageUnreadListCacheMgr == nil {
		return
	}
	if list == nil {
		list = make([]*entity.UserSystemMessageUnread, 0)
	}
	systemMessageUnreadListCacheMgr.FlushCache(userId, list)
}
