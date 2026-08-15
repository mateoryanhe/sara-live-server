package message

import (
	"context"
	"fmt"
	"time"

	"xr-game-server/core/cache"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity/message"
	"xr-game-server/module/upload"
)

const (
	privateMessageUnreadListPageSize = 50
	privateMessageUnreadListCacheMax = privateMessageUnreadListPageSize * 2
)

var privateMessageUnreadListCacheMgr = cache.NewCacheMgr()

func unreadListKey(userId uint64) string {
	return fmt.Sprintf("private_msg_unread_list:%d", userId)
}

func getPrivateMessageUnreadList(userId uint64) []*messagedto.AppPrivateMessageUnreadDetailItem {
	if userId == 0 {
		return nil
	}
	v := privateMessageUnreadListCacheMgr.GetData(unreadListKey(userId), func(ctx context.Context) (interface{}, error) {
		rows := messagedao.ListPrivateMessageUnreadWithLastMessageFromDBLimit(userId, privateMessageUnreadListCacheMax, 0)
		list := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(rows))
		for _, row := range rows {
			if row == nil {
				continue
			}
			item := toPrivateMessageUnreadDetailItemFromRow(row)
			fillPrivateMessageUnreadSenderInfo(item)
			list = append(list, item)
		}
		return list, nil
	})
	list, _ := v.([]*messagedto.AppPrivateMessageUnreadDetailItem)
	return list
}

func putPrivateMessageUnreadList(userId uint64, list []*messagedto.AppPrivateMessageUnreadDetailItem) {
	if userId == 0 {
		return
	}
	if list == nil {
		list = make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0)
	}
	privateMessageUnreadListCacheMgr.FlushCache(unreadListKey(userId), list)
}

func firstPagePrivateMessageUnreadList(userId uint64) []*messagedto.AppPrivateMessageUnreadDetailItem {
	list := getPrivateMessageUnreadList(userId)
	if len(list) > privateMessageUnreadListPageSize {
		return list[:privateMessageUnreadListPageSize]
	}
	return list
}

func prependPrivateMessageUnreadListCache(userId, senderId uint64, msg *entity.UserMessage, sessionId, unreadCount uint64) {
	if userId == 0 || senderId == 0 || msg == nil {
		return
	}
	list := getPrivateMessageUnreadList(userId)
	now := time.Now()
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = now
	}
	item := &messagedto.AppPrivateMessageUnreadDetailItem{
		SenderId: senderId, UnreadCount: unreadCount,
		UpdatedAt: formatMessageTime(now), LastMessage: toPrivateMessageItem(sessionId, msg),
	}
	fillPrivateMessageUnreadSenderInfo(item)

	newList := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(list)+1)
	newList = append(newList, item)
	for _, row := range list {
		if row != nil && row.SenderId != senderId {
			newList = append(newList, row)
		}
	}
	if len(newList) > privateMessageUnreadListCacheMax {
		newList = newList[:privateMessageUnreadListCacheMax]
	}
	putPrivateMessageUnreadList(userId, newList)
}

func updatePrivateMessageUnreadListCacheUnread(userId, senderId, unreadCount uint64) {
	list := getPrivateMessageUnreadList(userId)
	for _, item := range list {
		if item != nil && item.SenderId == senderId {
			item.UnreadCount = unreadCount
			item.UpdatedAt = formatMessageTime(time.Now())
			putPrivateMessageUnreadList(userId, list)
			return
		}
	}
}

func clearAllPrivateMessageUnreadListCacheUnread(userId uint64) {
	list := getPrivateMessageUnreadList(userId)
	now := formatMessageTime(time.Now())
	for _, item := range list {
		if item != nil {
			item.UnreadCount = 0
			item.UpdatedAt = now
		}
	}
	putPrivateMessageUnreadList(userId, list)
}

func removePrivateMessageUnreadListCacheSender(userId, senderId uint64) {
	list := getPrivateMessageUnreadList(userId)
	newList := make([]*messagedto.AppPrivateMessageUnreadDetailItem, 0, len(list))
	for _, item := range list {
		if item != nil && item.SenderId != senderId {
			newList = append(newList, item)
		}
	}
	putPrivateMessageUnreadList(userId, newList)
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
