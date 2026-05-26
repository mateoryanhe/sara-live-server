package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	privateMessagePageSize    = 40
	privateMessageCacheSize   = 200
	privateMessageCachedPages = 5
)

var receiverCacheMgr *cache.CacheMgr

// InitPrivateMessageDao 初始化私信缓存
func initPrivateMessageDao() {
	receiverCacheMgr = cache.NewCacheMgr()
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

	offset := (pageIndex - 1) * privateMessagePageSize
	return loadFromDB(receiverId, offset, privateMessagePageSize)
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
