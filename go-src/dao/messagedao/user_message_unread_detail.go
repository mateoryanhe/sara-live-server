package messagedao

import (
	"context"
	"time"

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
func GetUnreadDetailByReceiverSender(senderId, userId uint64) *entity.UserMessageUnreadDetail {

	id := entity.BuildUserMessageUnreadDetailId(userId, senderId)
	v := messageUnreadDetailCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.UserMessageUnreadDetail
		if err = g.Model(string(entity.TbUserMessageUnreadDetail)).Where("id = ?", id).Scan(&row); err != nil {
			return entity.NewUserMessageUnreadDetail(userId, senderId), err
		}
		if row == nil || row.ID == "" {
			return entity.NewUserMessageUnreadDetail(userId, senderId), nil
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
		list := GetCachedUnreadDetails(receiverId)
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
	if detail == nil || detail.ID == "" || messageUnreadDetailCacheMgr == nil || receiverUnreadDetailListCacheMgr == nil {
		return
	}
	//messageUnreadDetailCacheMgr.FlushCache(detail.ID, detail)

	list := GetCachedUnreadDetails(detail.UserId)
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
	receiverUnreadDetailListCacheMgr.FlushCache(detail.UserId, newList)
}

func GetCachedUnreadDetails(userId uint64) []*entity.UserMessageUnreadDetail {
	v := receiverUnreadDetailListCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		return loadUnreadDetailsFromDB(userId, 0, unreadDetailCacheSize), nil
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

func loadUnreadDetailsFromDB(userId uint64, offset, limit int) []*entity.UserMessageUnreadDetail {
	list := make([]*entity.UserMessageUnreadDetail, 0)
	_ = g.Model(string(entity.TbUserMessageUnreadDetail)).Ctx(context.Background()).
		Where("user_id = ?", userId).
		Order("updated_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list)
	return list
}

// PrivateMessageUnreadListRow App端私信未读列表查询行(未读明细 + 会话最后一条消息)
type PrivateMessageUnreadListRow struct {
	SenderId          uint64    `json:"sender_id"`
	UnreadCount       uint64    `json:"unread_count"`
	UpdatedAt         time.Time `json:"updated_at"`
	MessageId         uint64    `json:"message_id"`
	MessageSenderId   uint64    `json:"message_sender_id"`
	MessageReceiverId uint64    `json:"message_receiver_id"`
	MessageContent    string    `json:"message_content"`
	MessageCreatedAt  time.Time `json:"message_created_at"`
}

const listPrivateMessageUnreadWithLastMessageSQL = `
SELECT
  d.sender_id,
  d.unread_count,
  d.updated_at,
  lm.id AS message_id,
  lm.sender_id AS message_sender_id,
  lm.receiver_id AS message_receiver_id,
  lm.content AS message_content,
  lm.created_at AS message_created_at
FROM ` + string(entity.TbUserMessageUnreadDetail) + ` d
LEFT JOIN ` + string(entity.TbUserMessage) + ` lm ON lm.id = (
  SELECT m.id
  FROM ` + string(entity.TbUserMessage) + ` m
  INNER JOIN ` + string(entity.TbUserMessageSession) + ` s ON s.message_id = m.id
  WHERE s.session_id = d.id AND m.type = ?
  ORDER BY m.created_at DESC
  LIMIT 1
)
WHERE d.user_id = ?
ORDER BY d.updated_at DESC
LIMIT ? OFFSET ?`

// ListPrivateMessageUnreadWithLastMessageFromDB App端私信未读列表(直接查库,表关联查询最后一条消息)
func ListPrivateMessageUnreadWithLastMessageFromDB(userId uint64, pageIndex int) []*PrivateMessageUnreadListRow {
	list := make([]*PrivateMessageUnreadListRow, 0)
	if userId == 0 {
		return list
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * unreadDetailPageSize
	_ = g.DB().Ctx(context.Background()).Raw(
		listPrivateMessageUnreadWithLastMessageSQL,
		entity.UserMessageTypePrivate,
		userId,
		unreadDetailPageSize,
		offset,
	).Scan(&list)
	return list
}
