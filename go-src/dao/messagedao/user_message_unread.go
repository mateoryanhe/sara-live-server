package messagedao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/message"
)

var messageUnreadCacheMgr *cache.RowCache[*entity.UserMessageUnread]

// InitMessageUnreadDao 初始化用户消息未读缓存
func initMessageUnreadDao() {
	messageUnreadCacheMgr = cache.NewRowCache[*entity.UserMessageUnread]()
}

// PublishMessageUnread 原地修改未读汇总后刷新缓存.
func PublishMessageUnread(data *entity.UserMessageUnread) {
	if data == nil || data.ID == 0 || messageUnreadCacheMgr == nil {
		return
	}
	messageUnreadCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

// GetUnReadByUserId 获取用户消息未读,优先读缓存,缓存未命中再查库
func GetUnReadByUserId(userId uint64) *entity.UserMessageUnread {
	if userId == 0 || messageUnreadCacheMgr == nil {
		return nil
	}
	v := messageUnreadCacheMgr.MustGetRow(gctx.New(), userId, func(ctx context.Context) (*entity.UserMessageUnread, error) {
		var row *entity.UserMessageUnread
		err := g.Model(string(entity.TbUserMessageUnread)).Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&row)
		if err != nil {
			return entity.NewUserMessageUnread(userId), err
		}
		if row != nil && row.ID != 0 {
			return row, nil
		}
		return entity.NewUserMessageUnread(userId), nil
	})
	if v == nil || v.ID == 0 {
		return nil
	}
	return v
}
