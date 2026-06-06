package liveroomdao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var (
	// onlineCacheMgr 按复合ID(userId_roomId)缓存单条在线记录
	onlineCacheMgr *cache.CacheMgr
)

// InitLiveRoomOnlineDao 初始化在线玩家相关缓存
func initLiveRoomOnlineDao() {
	onlineCacheMgr = cache.NewCacheMgr()
}

// GetOnlineById 按复合ID获取在线记录(走缓存);数据库不存在则返回占位实例
func GetOnlineById(id string, userId, roomId uint64) *entity.LiveRoomOnline {
	if id == "" || onlineCacheMgr == nil {
		return nil
	}
	v := onlineCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var o *entity.LiveRoomOnline
		_ = g.Model(string(entity.TbLiveRoomOnline)).Where("id = ?", id).Scan(&o)
		if o == nil {
			return entity.NewLiveRoomOnline(userId, roomId), nil
		}
		return o, nil
	})
	if v == nil {
		return nil
	}
	o, _ := v.(*entity.LiveRoomOnline)
	return o
}

func GetAllOnline() []*entity.LiveRoomOnline {
	ret := make([]*entity.LiveRoomOnline, 0)
	g.Model(string(entity.TbLiveRoomOnline)).Where("status = ?", entity.LiveRoomOnlineStatusOnline).Scan(&ret)
	return ret
}
