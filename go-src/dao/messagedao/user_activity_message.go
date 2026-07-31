package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const UserActivityMessageListCacheASize = 50

var userActivityMessageListCacheAMgr *cache.CacheMgr

func initUserActivityMessageDao() {
	userActivityMessageListCacheAMgr = cache.NewCacheMgr()
}

// GetUserActivityMessageListCacheA 获取用户活动消息列表缓存A(前50条,按发布时间倒序)
func GetUserActivityMessageListCacheA(userId uint64) []*entity.UserActivityMessage {
	if userId == 0 || userActivityMessageListCacheAMgr == nil {
		return nil
	}
	v := userActivityMessageListCacheAMgr.GetData(userId, func(ctx context.Context) (interface{}, error) {
		list := make([]*entity.UserActivityMessage, 0)
		_ = g.Model(string(entity.TbUserActivityMessage)).Ctx(context.Background()).
			Where("user_id = ?", userId).
			Order("published_at desc").
			Limit(UserActivityMessageListCacheASize).
			Scan(&list)
		return list, nil
	})
	list, _ := v.([]*entity.UserActivityMessage)
	return list
}

// FlushUserActivityMessageListCacheA 刷新用户活动消息列表缓存A
func FlushUserActivityMessageListCacheA(userId uint64, list []*entity.UserActivityMessage) {
	if userId == 0 || userActivityMessageListCacheAMgr == nil {
		return
	}
	if list == nil {
		list = make([]*entity.UserActivityMessage, 0)
	}
	userActivityMessageListCacheAMgr.FlushCache(userId, list)
}

// RemoveUserActivityMessageByActivityMessageId 删除引用指定活动消息的用户记录并清空列表缓存A
func RemoveUserActivityMessageByActivityMessageId(activityMessageId uint64) error {
	if activityMessageId == 0 {
		return nil
	}
	_, err := g.Model(string(entity.TbUserActivityMessage)).Ctx(context.Background()).
		Where("activity_message_id = ?", activityMessageId).
		Delete()
	if err != nil {
		return err
	}
	ClearAllUserActivityMessageListCacheA()
	return nil
}

// ClearAllUserActivityMessageListCacheA 清空用户活动消息列表缓存A
func ClearAllUserActivityMessageListCacheA() {
	if userActivityMessageListCacheAMgr == nil || userActivityMessageListCacheAMgr.Cache == nil {
		return
	}
	_ = userActivityMessageListCacheAMgr.Cache.Clear(gctx.New())
}
