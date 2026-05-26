package messagedao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var messageUnreadDetailCacheMgr *cache.CacheMgr

func initMessageUnreadDetailDao() {
	messageUnreadDetailCacheMgr = cache.NewCacheMgr()
}

// GetUnreadDetailByReceiverSender 按接收者与发送者查询未读明细,优先读缓存,缓存未命中再查库
func GetUnreadDetailByReceiverSender(receiverId, senderId uint64) *entity.UserMessageUnreadDetail {
	if receiverId == 0 || senderId == 0 || messageUnreadDetailCacheMgr == nil {
		return nil
	}
	id := entity.BuildUserMessageUnreadDetailId(receiverId, senderId)
	v := messageUnreadDetailCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.UserMessageUnreadDetail
		if err = g.Model(string(entity.TbUserMessageUnreadDetail)).Where("id = ?", id).Scan(&row); err != nil {
			return nil, err
		}
		if row == nil || row.ID == "" {
			return entity.BuildUserMessageUnreadDetailId(receiverId, senderId), nil
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
