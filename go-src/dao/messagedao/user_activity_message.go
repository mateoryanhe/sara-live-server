package messagedao

import (
	"context"
	"sort"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
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

// GetUserActivityMessageListFromCacheA 仅从内存缓存读取用户活动消息列表,未命中不访问数据库
func GetUserActivityMessageListFromCacheA(userId uint64) ([]*entity.UserActivityMessage, bool) {
	if userId == 0 || userActivityMessageListCacheAMgr == nil {
		return nil, false
	}
	v := userActivityMessageListCacheAMgr.GetFromCache(userId)
	if v == nil {
		return nil, false
	}
	list, _ := v.([]*entity.UserActivityMessage)
	if list == nil {
		list = make([]*entity.UserActivityMessage, 0)
	}
	return list, true
}

// GetCachedUserActivityMessageUserIds 获取已有活动消息列表缓存的用户ID
func GetCachedUserActivityMessageUserIds() []uint64 {
	if userActivityMessageListCacheAMgr == nil || userActivityMessageListCacheAMgr.Cache == nil {
		return nil
	}
	keys := userActivityMessageListCacheAMgr.Cache.MustKeys(gctx.New())
	userIds := make([]uint64, 0, len(keys))
	for _, key := range keys {
		userId, ok := key.(uint64)
		if !ok || userId == 0 {
			continue
		}
		userIds = append(userIds, userId)
	}
	return userIds
}

// FlushUserActivityMessageListCacheA 刷新用户活动消息列表缓存A
func FlushUserActivityMessageListCacheA(userId uint64, list []*entity.UserActivityMessage) {
	if userId == 0 || userActivityMessageListCacheAMgr == nil {
		return
	}
	if list == nil {
		list = make([]*entity.UserActivityMessage, 0)
	}
	if len(list) > UserActivityMessageListCacheASize {
		list = list[:UserActivityMessageListCacheASize]
	}
	userActivityMessageListCacheAMgr.FlushCache(userId, list)
}

// PrependUserActivityMessageToListCacheA 向已有缓存头部插入活动消息,未命中缓存不访问数据库
func PrependUserActivityMessageToListCacheA(userId uint64, row *entity.UserActivityMessage) bool {
	if userId == 0 || row == nil || row.ActivityMessageId == 0 || userActivityMessageListCacheAMgr == nil {
		return false
	}
	list, ok := GetUserActivityMessageListFromCacheA(userId)
	if !ok {
		return false
	}
	for _, item := range list {
		if item != nil && item.ActivityMessageId == row.ActivityMessageId {
			return false
		}
	}
	newList := make([]*entity.UserActivityMessage, 0, len(list)+1)
	newList = append(newList, row)
	for _, item := range list {
		if item != nil && item.ActivityMessageId != row.ActivityMessageId {
			newList = append(newList, item)
		}
	}
	sortUserActivityMessageListByPublishedAtDesc(newList)
	if len(newList) > UserActivityMessageListCacheASize {
		newList = newList[:UserActivityMessageListCacheASize]
	}
	userActivityMessageListCacheAMgr.FlushCache(userId, newList)
	return true
}

// RemoveUserActivityMessageFromListCacheA 从已有缓存移除活动消息,未命中缓存不访问数据库
func RemoveUserActivityMessageFromListCacheA(userId uint64, activityMessageId uint64) bool {
	if userId == 0 || activityMessageId == 0 || userActivityMessageListCacheAMgr == nil {
		return false
	}
	list, ok := GetUserActivityMessageListFromCacheA(userId)
	if !ok {
		return false
	}
	newList := make([]*entity.UserActivityMessage, 0, len(list))
	found := false
	for _, item := range list {
		if item != nil && item.ActivityMessageId == activityMessageId {
			found = true
			continue
		}
		if item != nil {
			newList = append(newList, item)
		}
	}
	if !found {
		return false
	}
	userActivityMessageListCacheAMgr.FlushCache(userId, newList)
	return true
}

func sortUserActivityMessageListByPublishedAtDesc(list []*entity.UserActivityMessage) {
	sort.Slice(list, func(i, j int) bool {
		return userActivityMessagePublishedAt(list[i]).After(userActivityMessagePublishedAt(list[j]))
	})
}

func userActivityMessagePublishedAt(row *entity.UserActivityMessage) time.Time {
	if row == nil || row.PublishedAt == nil {
		return time.Time{}
	}
	return *row.PublishedAt
}

// RemoveUserActivityMessageByActivityMessageId 删除引用指定活动消息的用户记录
func RemoveUserActivityMessageByActivityMessageId(activityMessageId uint64) error {
	if activityMessageId == 0 {
		return nil
	}
	_, err := g.Model(string(entity.TbUserActivityMessage)).Ctx(context.Background()).
		Where("activity_message_id = ?", activityMessageId).
		Delete()
	return err
}

// ClearAllUserActivityMessageListCacheA 清空用户活动消息列表缓存A
func ClearAllUserActivityMessageListCacheA() {
	if userActivityMessageListCacheAMgr == nil || userActivityMessageListCacheAMgr.Cache == nil {
		return
	}
	_ = userActivityMessageListCacheAMgr.Cache.Clear(gctx.New())
}
