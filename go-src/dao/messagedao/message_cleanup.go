package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/message"
)

const userMessageMaxCount = 100000

const (
	deleteOldestMessageSessionsSQL = `
DELETE s FROM user_message_sessions s
INNER JOIN (
	SELECT id FROM user_messages ORDER BY created_at ASC, id ASC LIMIT ?
) t ON s.message_id = t.id`

	deleteOldestUserMessagesSQL = `
DELETE FROM user_messages WHERE id IN (
	SELECT id FROM (
		SELECT id FROM user_messages ORDER BY created_at ASC, id ASC LIMIT ?
	) t
)`
)

// CleanupExcessUserMessages 保留最新10万条消息,一次性删除更早数据及关联 session
func CleanupExcessUserMessages(ctx context.Context) (deletedMessages int, deletedSessions int) {
	total, err := g.Model(string(entity.TbUserMessage)).Ctx(ctx).Count()
	if err != nil || total <= userMessageMaxCount {
		return 0, 0
	}

	deleteCount := total - userMessageMaxCount

	sessionAffected, err := deleteOldestMessageSessionsRaw(ctx, deleteCount)
	if err != nil {
		g.Log().Errorf(ctx, "CleanupExcessUserMessages delete sessions err=%v count=%d", err, deleteCount)
		return 0, 0
	}

	messageAffected, err := deleteOldestUserMessagesRaw(ctx, deleteCount)
	if err != nil {
		g.Log().Errorf(ctx, "CleanupExcessUserMessages delete messages err=%v count=%d", err, deleteCount)
		return int(sessionAffected), 0
	}

	return int(messageAffected), int(sessionAffected)
}

func deleteOldestMessageSessionsRaw(ctx context.Context, limit int) (int64, error) {
	if limit <= 0 {
		return 0, nil
	}
	result, err := g.DB().Ctx(ctx).Exec(ctx, deleteOldestMessageSessionsSQL, limit)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func deleteOldestUserMessagesRaw(ctx context.Context, limit int) (int64, error) {
	if limit <= 0 {
		return 0, nil
	}
	result, err := g.DB().Ctx(ctx).Exec(ctx, deleteOldestUserMessagesSQL, limit)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
