package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const (
	deletePrivateMessageSessionsByTargetSQL = `
DELETE s FROM user_message_sessions s
INNER JOIN user_messages m ON m.id = s.message_id
WHERE m.type = ?
  AND ((m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?))`

	deletePrivateMessagesByTargetSQL = `
DELETE FROM user_messages
WHERE type = ?
  AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))`
)

// DeletePrivateConversationByTarget 一键删除与目标用户的私信(物理删除 session 与 message)
func DeletePrivateConversationByTarget(ctx context.Context, userId, targetId uint64) error {
	if userId == 0 || targetId == 0 {
		return nil
	}
	if _, err := DeletePrivateMessagesByTarget(ctx, userId, targetId); err != nil {
		return err
	}
	clearPrivateUnreadDetail(ctx, userId, targetId)
	return nil
}

// DeletePrivateMessagesByTarget 物理删除两人之间的全部私信及关联 session(无前置 id 查询)
func DeletePrivateMessagesByTarget(ctx context.Context, userId, targetId uint64) (deletedMessages int64, err error) {
	if userId == 0 || targetId == 0 {
		return 0, nil
	}
	msgType := entity.UserMessageTypePrivate
	args := []any{msgType, userId, targetId, targetId, userId}

	if _, err = g.DB().Ctx(ctx).Exec(ctx, deletePrivateMessageSessionsByTargetSQL, args...); err != nil {
		return 0, err
	}

	result, err := g.DB().Ctx(ctx).Exec(ctx, deletePrivateMessagesByTargetSQL, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func clearPrivateUnreadDetail(ctx context.Context, userId, targetId uint64) {
	detailId := entity.BuildUserMessageUnreadDetailId(userId, targetId)
	var detail *entity.UserMessageUnreadDetail
	if err := g.Model(string(entity.TbUserMessageUnreadDetail)).Ctx(ctx).
		Where("id = ?", detailId).
		Scan(&detail); err != nil || detail == nil || detail.ID == "" {
		return
	}
	if detail.UnreadCount == 0 {
		return
	}
	detail.ClearUnread(detail.UnreadCount)
}
