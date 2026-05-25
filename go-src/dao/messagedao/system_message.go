package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	systemMessagePageSize    = 40
	systemMessageCacheSize   = 200
	systemMessageCachedPages = 5
)

var systemReceiverCacheMgr *cache.CacheMgr

// InitSystemMessageDao 初始化系统消息缓存
func InitSystemMessageDao() {
	systemReceiverCacheMgr = cache.NewCacheMgr()
}

// ListSystemByReceiverId 按接收者分页查询系统消息
// 前5页(每页40条)走缓存,第6页及以后直接查库
func ListSystemByReceiverId(receiverId uint64, pageIndex int) []*entity.UserMessage {
	if receiverId == 0 || systemReceiverCacheMgr == nil {
		return make([]*entity.UserMessage, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}

	if pageIndex <= systemMessageCachedPages {
		list := getSystemCachedMessages(receiverId)
		start := (pageIndex - 1) * systemMessagePageSize
		if start >= len(list) {
			return make([]*entity.UserMessage, 0)
		}
		end := start + systemMessagePageSize
		if end > len(list) {
			end = len(list)
		}
		return list[start:end]
	}

	offset := (pageIndex - 1) * systemMessagePageSize
	return loadSystemFromDB(receiverId, offset, systemMessagePageSize)
}

// AddSystemToCache 新增系统消息后写入接收者缓存,超出200条时移除最旧数据
func AddSystemToCache(msg *entity.UserMessage) {
	if msg == nil || systemReceiverCacheMgr == nil || msg.ReceiverId == 0 {
		return
	}
	if msg.Type != entity.UserMessageTypeSystem {
		return
	}

	list := getSystemCachedMessages(msg.ReceiverId)
	systemReceiverCacheMgr.FlushCache(msg.ReceiverId, prependSystemMessage(list, msg))
}

func getSystemCachedMessages(receiverId uint64) []*entity.UserMessage {
	v := systemReceiverCacheMgr.GetData(receiverId, func(ctx context.Context) (value interface{}, err error) {
		return loadSystemFromDB(receiverId, 0, systemMessageCacheSize), nil
	})
	if v == nil {
		return make([]*entity.UserMessage, 0)
	}
	list, _ := v.([]*entity.UserMessage)
	if list == nil {
		return make([]*entity.UserMessage, 0)
	}
	return list
}

func prependSystemMessage(list []*entity.UserMessage, msg *entity.UserMessage) []*entity.UserMessage {
	newList := make([]*entity.UserMessage, 0, len(list)+1)
	newList = append(newList, msg)
	for _, item := range list {
		if item != nil && item.ID == msg.ID {
			continue
		}
		newList = append(newList, item)
	}
	if len(newList) > systemMessageCacheSize {
		newList = newList[:systemMessageCacheSize]
	}
	return newList
}

func loadSystemFromDB(receiverId uint64, offset, limit int) []*entity.UserMessage {
	list := make([]*entity.UserMessage, 0)
	_ = g.Model(string(entity.TbUserMessage)).Ctx(context.Background()).
		Where("receiver_id = ? AND type = ?", receiverId, entity.UserMessageTypeSystem).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list)
	return list
}
