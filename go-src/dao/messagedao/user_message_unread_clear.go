package messagedao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/message"
)

// ClearAllPrivateUnreadInDB 将指定用户全部私信未读数清零(直接更新数据库)
func ClearAllPrivateUnreadInDB(ctx context.Context, userId uint64) error {
	if userId == 0 {
		return nil
	}
	now := time.Now()

	_, err := g.Model(string(entity.TbUserMessageUnreadDetail)).Ctx(ctx).
		Where("user_id = ?", userId).
		Data(g.Map{
			string(entity.UserMessageUnreadDetailUnreadCount): 0,
			string(db.UpdatedAtName):                          now,
		}).Update()
	return err
}

// ClearAllPrivateUnreadCache 清零后刷新私信未读相关缓存
func ClearAllPrivateUnreadCache(userId uint64) {
	if userId == 0 {
		return
	}

	if messageUnreadDetailCacheMgr == nil {
		return
	}
	for _, item := range listCachedUnreadDetailsByUserId(userId) {
		if item == nil || item.UnreadCount == 0 {
			continue
		}
		item.UnreadCount = 0
		//messageUnreadDetailCacheMgr.FlushCache(item.ID, item)
	}
}
