package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	unreadDetailPageSize    = 40
	unreadDetailCacheSize   = 200
	unreadDetailCachedPages = 5
)

var messageUnreadDetailCacheMgr *cache.CacheMgr
var receiverUnreadDetailListCacheMgr *cache.CacheMgr

func initMessageUnreadDetailDao() {
	messageUnreadDetailCacheMgr = cache.NewCacheMgr()
	receiverUnreadDetailListCacheMgr = cache.NewCacheMgr()
}

// GetUnreadDetailByReceiverSender 按接收者与发送者查询未读明细,优先读缓存,缓存未命中再查库
func GetUnreadDetailByReceiverSender(receiverId, senderId uint64) *entity.UserMessageUnreadDetail {
	if receiverId == 0 || senderId == 0 || messageUnreadDetailCacheMgr == nil {
		return nil
	}
	id := entity.BuildUserMessageUnreadDetailId(receiverId, senderId)
	v := messageUnreadDetailCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.UserMessageUnreadDetail
		if err = g.Model(string(entity.TbUserMessageUnreadDetail)).Where("id = ?", id).Scan(&row); err != nil {
			return entity.NewUserMessageUnreadDetail(receiverId, senderId), err
		}
		if row == nil || row.ID == "" {
			return entity.NewUserMessageUnreadDetail(receiverId, senderId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.UserMessageUnreadDetail)
	if row == nil || row.ID == "" {
		return nil
	}
	return row
}

// ListUnreadDetailByReceiverId 按接收者分页查询未读明细
// 前200条(5页,每页40条)走缓存,第6页及以后直接查库
func ListUnreadDetailByReceiverId(receiverId uint64, pageIndex int) []*entity.UserMessageUnreadDetail {
	if receiverId == 0 {
		return make([]*entity.UserMessageUnreadDetail, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}

	if pageIndex <= unreadDetailCachedPages {
		list := getCachedUnreadDetails(receiverId)
		start := (pageIndex - 1) * unreadDetailPageSize
		if start >= len(list) {
			return make([]*entity.UserMessageUnreadDetail, 0)
		}
		end := start + unreadDetailPageSize
		if end > len(list) {
			end = len(list)
		}
		return list[start:end]
	}

	offset := (pageIndex - 1) * unreadDetailPageSize
	return loadUnreadDetailsFromDB(receiverId, offset, unreadDetailPageSize)
}

// UpsertUnreadDetailToListCache 未读明细变更后刷新单条缓存及接收者列表缓存
func UpsertUnreadDetailToListCache(detail *entity.UserMessageUnreadDetail) {
	if detail == nil || detail.ID == "" || detail.ReceiverId == 0 || messageUnreadDetailCacheMgr == nil || receiverUnreadDetailListCacheMgr == nil {
		return
	}
	//messageUnreadDetailCacheMgr.FlushCache(detail.ID, detail)

	list := getCachedUnreadDetails(detail.ReceiverId)
	newList := make([]*entity.UserMessageUnreadDetail, 0, len(list)+1)
	newList = append(newList, detail)
	for _, item := range list {
		if item == nil || item.ID == detail.ID {
			continue
		}
		newList = append(newList, item)
	}
	if len(newList) > unreadDetailCacheSize {
		newList = newList[:unreadDetailCacheSize]
	}
	receiverUnreadDetailListCacheMgr.FlushCache(detail.ReceiverId, newList)
}

func getCachedUnreadDetails(receiverId uint64) []*entity.UserMessageUnreadDetail {
	v := receiverUnreadDetailListCacheMgr.GetData(receiverId, func(ctx context.Context) (value interface{}, err error) {
		return loadUnreadDetailsFromDB(receiverId, 0, unreadDetailCacheSize), nil
	})
	if v == nil {
		return make([]*entity.UserMessageUnreadDetail, 0)
	}
	list, _ := v.([]*entity.UserMessageUnreadDetail)
	if list == nil {
		return make([]*entity.UserMessageUnreadDetail, 0)
	}
	return list
}

func loadUnreadDetailsFromDB(receiverId uint64, offset, limit int) []*entity.UserMessageUnreadDetail {
	list := make([]*entity.UserMessageUnreadDetail, 0)
	_ = g.Model(string(entity.TbUserMessageUnreadDetail)).Ctx(context.Background()).
		Where("receiver_id = ?", receiverId).
		Order("updated_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list)
	return list
}
