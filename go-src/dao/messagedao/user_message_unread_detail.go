package messagedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
)

const unreadDetailPageSize = 40

var messageUnreadDetailCacheMgr *cache.CacheMgr

func initMessageUnreadDetailDao() {
	messageUnreadDetailCacheMgr = cache.NewCacheMgr()
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

// FlushUnreadDetailCache 刷新单条未读明细缓存
func FlushUnreadDetailCache(detail *entity.UserMessageUnreadDetail) {
	if detail == nil || detail.ID == "" || messageUnreadDetailCacheMgr == nil {
		return
	}
	messageUnreadDetailCacheMgr.FlushCache(detail.ID, detail)
}

func listCachedUnreadDetailsByUserId(userId uint64) []*entity.UserMessageUnreadDetail {
	if messageUnreadDetailCacheMgr == nil || userId == 0 {
		return nil
	}
	ctx := gctx.New()
	keys, err := messageUnreadDetailCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return nil
	}

	list := make([]*entity.UserMessageUnreadDetail, 0)
	for _, key := range keys {
		value, err := messageUnreadDetailCacheMgr.Cache.Get(ctx, key)
		if err != nil || value == nil || value.IsNil() {
			continue
		}
		item, ok := value.Val().(*entity.UserMessageUnreadDetail)
		if !ok || item == nil || item.UserId != userId {
			continue
		}
		list = append(list, item)
	}
	return list
}

// PrivateMessageUnreadListRow App端私信未读列表查询行(未读明细 + 会话最后一条消息)
type PrivateMessageUnreadListRow struct {
	SenderId          uint64    `json:"sender_id"`
	UnreadCount       uint64    `json:"unread_count"`
	UpdatedAt         time.Time `json:"updated_at"`
	SessionRowId      uint64    `json:"session_row_id"`
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
  ls.id AS session_row_id,
  lm.id AS message_id,
  lm.sender_id AS message_sender_id,
  lm.receiver_id AS message_receiver_id,
  lm.content AS message_content,
  lm.created_at AS message_created_at
FROM ` + string(entity.TbUserMessageUnreadDetail) + ` d
LEFT JOIN ` + string(entity.TbUserMessageSession) + ` ls ON ls.id = (
  SELECT s.id
  FROM ` + string(entity.TbUserMessageSession) + ` s
  INNER JOIN ` + string(entity.TbUserMessage) + ` m ON s.message_id = m.id
  WHERE s.session_id = d.id AND s.is_deleted = 0 AND m.is_deleted = 0 AND m.sender_id > 0
  ORDER BY m.created_at DESC
  LIMIT 1
)
LEFT JOIN ` + string(entity.TbUserMessage) + ` lm ON lm.id = ls.message_id AND lm.is_deleted = 0
WHERE d.user_id = ? AND d.mutual_chat = ?
ORDER BY d.updated_at DESC
LIMIT ? OFFSET ?`

// ListPrivateMessageUnreadWithLastMessageFromDB App端私信未读列表(直接查库,表关联查询最后一条消息)
func ListPrivateMessageUnreadWithLastMessageFromDB(userId uint64, pageIndex int) []*PrivateMessageUnreadListRow {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	offset := (pageIndex - 1) * unreadDetailPageSize
	return ListPrivateMessageUnreadWithLastMessageFromDBLimit(userId, unreadDetailPageSize, offset)
}

// ListPrivateMessageUnreadWithLastMessageFromDBLimit App端私信未读列表(直接查库,limit/offset)
func ListPrivateMessageUnreadWithLastMessageFromDBLimit(userId uint64, limit, offset int) []*PrivateMessageUnreadListRow {
	list := make([]*PrivateMessageUnreadListRow, 0)
	if userId == 0 || limit <= 0 {
		return list
	}
	if offset < 0 {
		offset = 0
	}
	_ = g.DB().Ctx(context.Background()).Raw(
		listPrivateMessageUnreadWithLastMessageSQL,
		userId,
		entity.UserMessageUnreadMutualChatYes,
		limit,
		offset,
	).Scan(&list)
	return list
}

// MarkMutualPrivateChat 标记双方未读明细为已互发私信
func MarkMutualPrivateChat(userA, userB uint64) {
	if userA == 0 || userB == 0 || userA == userB {
		return
	}
	markOneMutualPrivateChat(userA, userB)
	markOneMutualPrivateChat(userB, userA)
}

func markOneMutualPrivateChat(receiverId, senderId uint64) {
	detail := GetUnreadDetailByReceiverSender(senderId, receiverId)
	if detail == nil || detail.MutualChat == entity.UserMessageUnreadMutualChatYes {
		return
	}
	detail.SetMutualChat(entity.UserMessageUnreadMutualChatYes)
	FlushUnreadDetailCache(detail)
}
