package messagedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const privateMessagePageSize = 40

func initPrivateMessageDao() {}

// ListByReceiverAndSender 按会话ID查询私信(关联 user_messages + user_message_sessions)
// lastCreatedAt=0 从最新消息开始; >0 拉取创建时间更早的消息
func ListByReceiverAndSender(sessionId string, lastCreatedAt int64, pageSize int) ([]*entity.UserMessage, bool) {
	if sessionId == "" {
		return make([]*entity.UserMessage, 0), false
	}
	if lastCreatedAt < 0 {
		lastCreatedAt = 0
	}
	if pageSize <= 0 {
		pageSize = privateMessagePageSize
	}

	limit := pageSize + 1
	model := g.Model(string(entity.TbUserMessage)+" m").Ctx(context.Background()).
		Fields("m.*").
		InnerJoin(string(entity.TbUserMessageSession)+" s", "s.message_id = m.id").
		Where("s.session_id = ? AND m.type = ?", sessionId, entity.UserMessageTypePrivate)
	if lastCreatedAt > 0 {
		model = model.Where("m.created_at < ?", time.UnixMilli(lastCreatedAt))
	}

	list := make([]*entity.UserMessage, 0, limit)
	if err := model.Order("m.created_at desc").Limit(limit).Scan(&list); err != nil {
		return make([]*entity.UserMessage, 0), false
	}

	hasMore := len(list) > pageSize
	if hasMore {
		list = list[:pageSize]
	}
	return list, hasMore
}
