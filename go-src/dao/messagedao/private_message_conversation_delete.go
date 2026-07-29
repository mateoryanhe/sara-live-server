package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const deletePrivateMessageSessionsByUserSQL = `
DELETE FROM user_message_sessions
WHERE session_id = ?`

// DeletePrivateConversationByTarget 一键删除与目标用户的私信(物理删除当前用户侧 session 索引)
func DeletePrivateConversationByTarget(ctx context.Context, userId, targetId uint64) error {
	if userId == 0 || targetId == 0 {
		return nil
	}
	_, err := g.DB().Ctx(ctx).Exec(
		ctx,
		deletePrivateMessageSessionsByUserSQL,
		entity.BuildUserMessageSessionId(userId, targetId),
	)
	return err
}
