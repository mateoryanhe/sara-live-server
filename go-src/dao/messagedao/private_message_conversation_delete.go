package messagedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const softDeletePrivateMessageSessionsByUserSQL = `
UPDATE user_message_sessions
SET is_deleted = 1, deleted_at = ?, updated_at = ?
WHERE session_id = ? AND is_deleted = 0`

// DeletePrivateConversationByTarget 一键删除与目标用户的私信(仅软删除当前用户侧 session)
func DeletePrivateConversationByTarget(ctx context.Context, userId, targetId uint64) error {
	if userId == 0 || targetId == 0 {
		return nil
	}
	_, err := SoftDeletePrivateMessageSessionsByUser(ctx, userId, targetId)
	return err
}

// SoftDeletePrivateMessageSessionsByUser 软删除当前用户视角下与目标用户的全部私信 session
func SoftDeletePrivateMessageSessionsByUser(ctx context.Context, userId, targetId uint64) (affected int64, err error) {
	if userId == 0 || targetId == 0 {
		return 0, nil
	}
	now := time.Now()
	result, err := g.DB().Ctx(ctx).Exec(
		ctx,
		softDeletePrivateMessageSessionsByUserSQL,
		now,
		now,
		entity.BuildUserMessageSessionId(userId, targetId),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
