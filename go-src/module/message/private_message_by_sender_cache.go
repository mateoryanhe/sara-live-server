package message

import (
	"context"
	"fmt"

	"xr-game-server/core/cache"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dto/messagedto"
	"xr-game-server/entity"
)

const (
	privateMessageBySenderPageSize = 40
	privateMessageBySenderCacheMax = privateMessageBySenderPageSize * 2
)

type privateMessageBySenderCacheData struct {
	Messages []*messagedto.AppPrivateMessageItem
	HasMore  bool
}

var privateMessageBySenderCacheMgr = cache.NewCacheMgr()

func bySenderCacheKey(userId, targetId uint64) string {
	return fmt.Sprintf("private_msg_by_sender:%d:%d", userId, targetId)
}

func getPrivateMessageBySender(userId, targetId uint64) *privateMessageBySenderCacheData {
	if userId == 0 || targetId == 0 {
		return &privateMessageBySenderCacheData{}
	}
	v := privateMessageBySenderCacheMgr.GetData(bySenderCacheKey(userId, targetId), func(ctx context.Context) (interface{}, error) {
		sessionId := entity.BuildUserMessageSessionId(userId, targetId)
		rows, hasMore := messagedao.ListByReceiverAndSender(sessionId, 0, privateMessageBySenderCacheMax)
		messages := make([]*messagedto.AppPrivateMessageItem, 0, len(rows))
		for _, row := range rows {
			if row == nil || row.Message == nil {
				continue
			}
			messages = append(messages, toPrivateMessageItem(row.SessionId, row.Message))
		}
		return &privateMessageBySenderCacheData{Messages: messages, HasMore: hasMore}, nil
	})
	data, _ := v.(*privateMessageBySenderCacheData)
	if data == nil {
		return &privateMessageBySenderCacheData{}
	}
	return data
}

func putPrivateMessageBySender(userId, targetId uint64, data *privateMessageBySenderCacheData) {
	if userId == 0 || targetId == 0 {
		return
	}
	if data == nil {
		data = &privateMessageBySenderCacheData{}
	}
	if data.Messages == nil {
		data.Messages = make([]*messagedto.AppPrivateMessageItem, 0)
	}
	privateMessageBySenderCacheMgr.FlushCache(bySenderCacheKey(userId, targetId), data)
}

func firstPagePrivateMessageBySender(userId, targetId uint64, pageSize int) ([]*messagedto.AppPrivateMessageItem, bool) {
	if pageSize <= 0 {
		pageSize = privateMessageBySenderPageSize
	}
	data := getPrivateMessageBySender(userId, targetId)
	list := data.Messages
	if len(list) > pageSize {
		list = list[:pageSize]
	}
	return list, len(data.Messages) > pageSize || data.HasMore
}

func prependPrivateMessageBySenderCache(userId, targetId, sessionId uint64, msg *entity.UserMessage) {
	if userId == 0 || targetId == 0 || msg == nil {
		return
	}
	data := getPrivateMessageBySender(userId, targetId)
	item := toPrivateMessageItem(sessionId, msg)
	newList := make([]*messagedto.AppPrivateMessageItem, 0, len(data.Messages)+1)
	newList = append(newList, item)
	for _, row := range data.Messages {
		if row != nil && row.Id != msg.ID {
			newList = append(newList, row)
		}
	}
	if len(newList) > privateMessageBySenderCacheMax {
		newList = newList[:privateMessageBySenderCacheMax]
	}
	putPrivateMessageBySender(userId, targetId, &privateMessageBySenderCacheData{
		Messages: newList,
		HasMore:  data.HasMore,
	})
}

func clearPrivateMessageBySenderCache(userId, targetId uint64) {
	putPrivateMessageBySender(userId, targetId, &privateMessageBySenderCacheData{})
}
