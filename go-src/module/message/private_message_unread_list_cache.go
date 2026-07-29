package message

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

const (
	privateMessageUnreadListPageSize       = 30
	privateMessageUnreadListCachePageCount = 2
	privateMessageUnreadListCacheMaxSize   = privateMessageUnreadListPageSize * privateMessageUnreadListCachePageCount
)

var privateMessageUnreadListCacheMgr *cache.CacheMgr

func initPrivateMessageUnreadListCache() {
	privateMessageUnreadListCacheMgr = cache.NewCacheMgr()
}

func privateMessageUnreadListCacheKey(userId uint64) string {
	return fmt.Sprintf("private_msg_unread_list:%d", userId)
}

func getPrivateMessageUnreadListCache(userId uint64) ([]*messagedto.AppPrivateMessageUnreadDetailItem, bool) {
	if userId == 0 || privateMessageUnreadListCacheMgr == nil {
		return nil, false
	}
	val, err := privateMessageUnreadListCacheMgr.Cache.Get(gctx.New(), privateMessageUnreadListCacheKey(userId))
	if err != nil || val == nil || val.IsNil() {
		return nil, false
	}
	list, ok := val.Val().([]*messagedto.AppPrivateMessageUnreadDetailItem)
	return list, ok
}

func setPrivateMessageUnreadListCache(userId uint64, list []*messagedto.AppPrivateMessageUnreadDetailItem) {
	if userId == 0 || privateMessageUnreadListCacheMgr == nil {
		return
	}
	privateMessageUnreadListCacheMgr.FlushCache(privateMessageUnreadListCacheKey(userId), list)
}

func pagePrivateMessageUnreadList(list []*messagedto.AppPrivateMessageUnreadDetailItem, pageIndex int) []*messagedto.AppPrivateMessageUnreadDetailItem {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	start := (pageIndex - 1) * privateMessageUnreadListPageSize
	if start >= len(list) {
		return make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0)
	}
	end := start + privateMessageUnreadListPageSize
	if end > len(list) {
		end = len(list)
	}
	ret := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, end-start)
	for _, item := range list[start:end] {
		if item == nil {
			continue
		}
		ret = append(ret, clonePrivateMessageUnreadDetailItem(item))
	}
	return ret
}

func loadPrivateMessageUnreadListCache(userId uint64) []*messagedto.AppPrivateMessageUnreadDetailItem {
	rows := messagedao.ListPrivateMessageUnreadWithLastMessageFromDBLimit(
		userId, privateMessageUnreadListCacheMaxSize, 0,
	)
	list := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := toPrivateMessageUnreadDetailItemFromRow(row)
		fillPrivateMessageUnreadSenderInfo(item)
		list = append(list, item)
	}
	setPrivateMessageUnreadListCache(userId, list)
	return list
}

func prependPrivateMessageUnreadListCache(
	receiverId, senderId uint64,
	msg *entity.UserMessage,
	sessionId uint64,
	unreadCount uint64,
) {
	if receiverId == 0 || senderId == 0 || msg == nil {
		return
	}
	list, _ := getPrivateMessageUnreadListCache(receiverId)
	now := time.Now()
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = now
	}
	item := &messagedto.AppPrivateMessageUnreadDetailItem{
		SenderId:    senderId,
		UnreadCount: unreadCount,
		UpdatedAt:   formatMessageTime(now),
		LastMessage: toPrivateMessageItem(sessionId, msg),
	}
	fillPrivateMessageUnreadSenderInfo(item)

	filtered := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(list)+1)
	filtered = append(filtered, item)
	for _, row := range list {
		if row == nil || row.SenderId == senderId {
			continue
		}
		filtered = append(filtered, row)
	}
	if len(filtered) > privateMessageUnreadListCacheMaxSize {
		filtered = filtered[:privateMessageUnreadListCacheMaxSize]
	}
	setPrivateMessageUnreadListCache(receiverId, filtered)
}

func updatePrivateMessageUnreadListCacheUnread(userId, senderId uint64, unreadCount uint64) {
	if userId == 0 || senderId == 0 {
		return
	}
	list, ok := getPrivateMessageUnreadListCache(userId)
	if !ok || len(list) == 0 {
		return
	}
	for _, item := range list {
		if item != nil && item.SenderId == senderId {
			item.UnreadCount = unreadCount
			item.UpdatedAt = formatMessageTime(time.Now())
			setPrivateMessageUnreadListCache(userId, list)
			return
		}
	}
}

func clearAllPrivateMessageUnreadListCacheUnread(userId uint64) {
	if userId == 0 {
		return
	}
	list, ok := getPrivateMessageUnreadListCache(userId)
	if !ok {
		list = loadPrivateMessageUnreadListCache(userId)
	}
	if len(list) == 0 {
		return
	}
	now := formatMessageTime(time.Now())
	for _, item := range list {
		if item == nil {
			continue
		}
		item.UnreadCount = 0
		item.UpdatedAt = now
	}
	setPrivateMessageUnreadListCache(userId, list)
}

func removePrivateMessageUnreadListCacheSender(userId, senderId uint64) {
	if userId == 0 || senderId == 0 {
		return
	}
	list, ok := getPrivateMessageUnreadListCache(userId)
	if !ok || len(list) == 0 {
		return
	}
	filtered := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(list))
	for _, item := range list {
		if item == nil || item.SenderId == senderId {
			continue
		}
		filtered = append(filtered, item)
	}
	if len(filtered) == 0 {
		removePrivateMessageUnreadListCache(userId)
		return
	}
	setPrivateMessageUnreadListCache(userId, filtered)
}

func removePrivateMessageUnreadListCache(userId uint64) {
	if userId == 0 || privateMessageUnreadListCacheMgr == nil {
		return
	}
	_, _ = privateMessageUnreadListCacheMgr.Cache.Remove(gctx.New(), privateMessageUnreadListCacheKey(userId))
}

func fillPrivateMessageUnreadSenderInfo(item *messagedto.AppPrivateMessageUnreadDetailItem) {
	if item == nil {
		return
	}
	if sender := userinfodao.GetUserInfoByUserId(item.SenderId); sender != nil {
		item.SenderName = sender.Nickname
		item.SenderAvatar = upload.ResolveAvatarUrlForUser(item.SenderId, sender.Avatar)
	}
	if item.LastMessage != nil {
		if sender := userinfodao.GetUserInfoByUserId(item.LastMessage.SenderId); sender != nil {
			item.LastMessage.SenderName = sender.Nickname
			item.LastMessage.SenderAvatar = upload.ResolveAvatarUrlForUser(item.LastMessage.SenderId, sender.Avatar)
		}
	}
}

func clonePrivateMessageUnreadDetailItem(item *messagedto.AppPrivateMessageUnreadDetailItem) *messagedto.AppPrivateMessageUnreadDetailItem {
	if item == nil {
		return nil
	}
	ret := *item
	if item.LastMessage != nil {
		lastMessage := *item.LastMessage
		ret.LastMessage = &lastMessage
	}
	fillPrivateMessageUnreadSenderInfo(&ret)
	return &ret
}
