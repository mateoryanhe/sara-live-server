package messagedao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
)

const SystemMessageUnreadListCacheMax = 100

var systemMessageUnreadListCacheMgr *cache.ListCache[*entity.UserSystemMessageUnread]

func initSystemMessageUnreadDao() {
	systemMessageUnreadListCacheMgr = cache.NewListCache[*entity.UserSystemMessageUnread]()
}

// GetSystemMessageUnreadListCache 获取用户系统消息未读列表(优先读缓存,未命中查库)
func GetSystemMessageUnreadListCache(userId uint64) []*entity.UserSystemMessageUnread {
	if userId == 0 || systemMessageUnreadListCacheMgr == nil {
		return nil
	}
	v := systemMessageUnreadListCacheMgr.MustGetList(gctx.New(), userId, func(ctx context.Context) ([]*entity.UserSystemMessageUnread, error) {
		list := make([]*entity.UserSystemMessageUnread, 0)
		_ = g.Model(string(entity.TbUserSystemMessageUnread)).Ctx(context.Background()).
			Where("user_id = ? AND unread_count > 0", userId).
			Order("updated_at desc").
			Limit(SystemMessageUnreadListCacheMax).
			Scan(&list)
		return list, nil
	})
	return v
}

// FlushSystemMessageUnreadListCache 刷新用户系统消息未读列表缓存
func FlushSystemMessageUnreadListCache(userId uint64, list []*entity.UserSystemMessageUnread) {
	if userId == 0 || systemMessageUnreadListCacheMgr == nil {
		return
	}
	if list == nil {
		list = make([]*entity.UserSystemMessageUnread, 0)
	}
	systemMessageUnreadListCacheMgr.PublishList(gctx.New(), userId, list)
}
