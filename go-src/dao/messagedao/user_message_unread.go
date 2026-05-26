package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var messageUnreadCacheMgr *cache.CacheMgr

// InitMessageUnreadDao 初始化用户消息未读缓存
func initMessageUnreadDao() {
	messageUnreadCacheMgr = cache.NewCacheMgr()
}

// GetUnReadByUserId 获取用户消息未读,优先读缓存,缓存未命中再查库
func GetUnReadByUserId(userId uint64) *entity.UserMessageUnread {
	if userId == 0 || messageUnreadCacheMgr == nil {
		return nil
	}
	v := messageUnreadCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.UserMessageUnread
		err = g.Model(string(entity.TbUserMessageUnread)).Where(g.Map{
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
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.UserMessageUnread)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}
