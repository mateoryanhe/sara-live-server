package messagedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const privateMessagePageSize = 40

// PrivateMessageListRow 私信列表行(含当前用户会话索引ID)
type PrivateMessageListRow struct {
	SessionId uint64              `json:"sessionId"`
	Message   *entity.UserMessage `json:"message"`
}

type privateMessageListScanRow struct {
	entity.UserMessage
	SessionId uint64 `json:"session_id"`
}

func initPrivateMessageDao() {}

// ListByReceiverAndSender 按会话ID查询私信(关联 user_messages + user_message_sessions)
// lastCreatedAt=0 从最新消息开始; >0 拉取创建时间更早的消息
func ListByReceiverAndSender(sessionId string, lastCreatedAt int64, pageSize int) ([]*PrivateMessageListRow, bool) {
	if sessionId == "" {
		return make([]*PrivateMessageListRow, 0), false
	}
	if lastCreatedAt < 0 {
		lastCreatedAt = 0
	}
	if pageSize <= 0 {
		pageSize = privateMessagePageSize
	}

	limit := pageSize + 1
	model := g.Model(string(entity.TbUserMessage)+" m").Ctx(context.Background()).
		Fields("m.*, s.id AS session_id").
		InnerJoin(string(entity.TbUserMessageSession)+" s", "s.message_id = m.id").
		Where("s.session_id = ? AND s.is_deleted = 0 AND m.is_deleted = 0 AND m.type = ?", sessionId, entity.UserMessageTypePrivate)
	if lastCreatedAt > 0 {
		model = model.Where("m.created_at < ?", time.UnixMilli(lastCreatedAt))
	}

	scanList := make([]*privateMessageListScanRow, 0, limit)
	if err := model.Order("m.created_at desc").Limit(limit).Scan(&scanList); err != nil {
		return make([]*PrivateMessageListRow, 0), false
	}

	hasMore := len(scanList) > pageSize
	if hasMore {
		scanList = scanList[:pageSize]
	}

	list := make([]*PrivateMessageListRow, 0, len(scanList))
	for _, row := range scanList {
		if row == nil {
			continue
		}
		msg := row.UserMessage
		list = append(list, &PrivateMessageListRow{
			SessionId: row.SessionId,
			Message:   &msg,
		})
	}
	return list, hasMore
}
