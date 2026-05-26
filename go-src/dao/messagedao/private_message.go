package messagedao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	privateMessagePageSize         = 40
	privateMessageCacheSize        = 200
	privateMessageCachedPages      = 5
	privateMessageBySenderMaxCount = 200
)

var receiverCacheMgr *cache.CacheMgr
var senderReceiverMessageCacheMgr *cache.CacheMgr

// InitPrivateMessageDao 初始化私信缓存
func initPrivateMessageDao() {
	receiverCacheMgr = cache.NewCacheMgr()
	senderReceiverMessageCacheMgr = cache.NewCacheMgr()
}

func buildSenderReceiverMessageCacheKey(receiverId, senderId uint64) string {
	return fmt.Sprintf("pm_%d_%d", receiverId, senderId)
}

// ListByReceiverAndSender 按接收者与发送者分页查询私信,最多缓存200条
// 优先读缓存,缓存未命中再查库;总条数超过200条时返回空列表
func ListByReceiverAndSender(receiverId, senderId uint64, pageIndex, pageSize int) []*entity.UserMessage {
	if receiverId == 0 || senderId == 0 || senderReceiverMessageCacheMgr == nil {
		return make([]*entity.UserMessage, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = privateMessagePageSize
	}
	if pageSize > privateMessageBySenderMaxCount {
		pageSize = privateMessageBySenderMaxCount
	}

	list := getCachedMessagesByReceiverSender(receiverId, senderId)
	start := (pageIndex - 1) * pageSize
	if start >= len(list) || start >= privateMessageBySenderMaxCount {
		return make([]*entity.UserMessage, 0)
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	if end > privateMessageBySenderMaxCount {
		end = privateMessageBySenderMaxCount
	}
	return list[start:end]
}

func getCachedMessagesByReceiverSender(receiverId, senderId uint64) []*entity.UserMessage {
	cacheKey := buildSenderReceiverMessageCacheKey(receiverId, senderId)
	v := senderReceiverMessageCacheMgr.GetData(cacheKey, func(ctx context.Context) (value interface{}, err error) {
		list := loadFromDBByReceiverSender(receiverId, senderId, 0, privateMessageBySenderMaxCount+1)
		if len(list) > privateMessageBySenderMaxCount {
			return make([]*entity.UserMessage, 0), nil
		}
		return list, nil
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

// InvalidateSenderReceiverMessageCache 清除指定会话私信缓存
func InvalidateSenderReceiverMessageCache(receiverId, senderId uint64) {
	if receiverId == 0 || senderId == 0 || senderReceiverMessageCacheMgr == nil {
		return
	}
	senderReceiverMessageCacheMgr.FlushCache(buildSenderReceiverMessageCacheKey(receiverId, senderId), nil)
}

// ListByReceiverId 按接收者分页查询私信,receiverId=用户ID
// 前5页(每页40条)走缓存,第6页及以后直接查库
func ListByReceiverId(receiverId uint64, pageIndex int) []*entity.UserMessage {
	if receiverId == 0 {
		return make([]*entity.UserMessage, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}

	if pageIndex <= privateMessageCachedPages {
		list := getCachedMessages(receiverId)
		start := (pageIndex - 1) * privateMessagePageSize
		if start >= len(list) {
			return make([]*entity.UserMessage, 0)
		}
		end := start + privateMessagePageSize
		if end > len(list) {
			end = len(list)
		}
		return list[start:end]
	}
	return make([]*entity.UserMessage, 0)
	//offset := (pageIndex - 1) * privateMessagePageSize
	//return loadFromDB(receiverId, offset, privateMessagePageSize)
}

// AddToCache 新增私信后写入接收者缓存,超出200条时移除最旧数据
func AddToCache(msg *entity.UserMessage, userId uint64) {
	if msg == nil || receiverCacheMgr == nil || msg.ReceiverId == 0 {
		return
	}
	if msg.Type != entity.UserMessageTypePrivate {
		return
	}

	list := getCachedMessages(userId)
	newList := make([]*entity.UserMessage, 0, len(list)+1)
	newList = append(newList, msg)
	for _, item := range list {
		if item != nil && item.ID == msg.ID {
			continue
		}
		newList = append(newList, item)
	}
	if len(newList) > privateMessageCacheSize {
		newList = newList[:privateMessageCacheSize]
	}
	receiverCacheMgr.FlushCache(userId, newList)
}

func getCachedMessages(receiverId uint64) []*entity.UserMessage {
	v := receiverCacheMgr.GetData(receiverId, func(ctx context.Context) (value interface{}, err error) {
		return loadFromDB(receiverId, 0, privateMessageCacheSize), nil
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

func loadFromDB(receiverId uint64, offset, limit int) []*entity.UserMessage {
	list := make([]*entity.UserMessage, 0)
	_ = g.Model(string(entity.TbUserMessage)).Ctx(context.Background()).
		Where("receiver_id = ? AND type = ?", receiverId, entity.UserMessageTypePrivate).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list)
	return list
}

func loadFromDBByReceiverSender(receiverId, senderId uint64, offset, limit int) []*entity.UserMessage {
	list := make([]*entity.UserMessage, 0)
	_ = g.Model(string(entity.TbUserMessage)).Ctx(context.Background()).
		Where("receiver_id = ? AND sender_id = ? AND type = ?", receiverId, senderId, entity.UserMessageTypePrivate).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list)
	return list
}
