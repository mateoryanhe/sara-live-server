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
	// roomOnlineCacheMgr 按 roomId 缓存在线玩家列表 []*entity.LiveRoomOnline
	roomOnlineCacheMgr *cache.CacheMgr
)

// InitLiveRoomOnlineDao 初始化在线玩家相关缓存
func InitLiveRoomOnlineDao() {
	onlineCacheMgr = cache.NewCacheMgr()
	roomOnlineCacheMgr = cache.NewCacheMgr()
}

// GetOnlineById 按复合ID获取在线记录(走缓存)
func GetOnlineById(id string) *entity.LiveRoomOnline {
	v := onlineCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var o *entity.LiveRoomOnline
		_ = g.Model(string(entity.TbLiveRoomOnline)).Where("id = ?", id).Scan(&o)
		return o, nil
	})
	if v == nil {
		return nil
	}
	o, _ := v.(*entity.LiveRoomOnline)
	return o
}

// GetOnlinesByRoom 按 roomId 获取"在线"玩家列表(走缓存,仅 Status == Online)
func GetOnlinesByRoom(roomId uint64) []*entity.LiveRoomOnline {
	v := roomOnlineCacheMgr.GetData(roomId, func(ctx context.Context) (value interface{}, err error) {
		list := make([]*entity.LiveRoomOnline, 0)
		_ = g.Model(string(entity.TbLiveRoomOnline)).
			Where("room_id = ? AND status = ?", roomId, entity.LiveRoomOnlineStatusOnline).
			Scan(&list)
		return list, nil
	})
	if v == nil {
		return nil
	}
	list, _ := v.([]*entity.LiveRoomOnline)
	return list
}

// AddOnlineToRoomCache 将"在线"记录写入房间在线列表缓存(去重)
func AddOnlineToRoomCache(o *entity.LiveRoomOnline) {
	if o == nil {
		return
	}
	onlineCacheMgr.FlushCache(o.ID, o)
	list := GetOnlinesByRoom(o.RoomId)
	for _, item := range list {
		if item.ID == o.ID {
			return
		}
	}
	list = append(list, o)
	roomOnlineCacheMgr.FlushCache(o.RoomId, list)
}

// RemoveOnlineFromRoomCache 将记录从房间在线列表缓存中剔除(单条缓存保留,便于下次加入复用)
func RemoveOnlineFromRoomCache(o *entity.LiveRoomOnline) {
	if o == nil || roomOnlineCacheMgr == nil {
		return
	}
	list := GetOnlinesByRoom(o.RoomId)
	newList := make([]*entity.LiveRoomOnline, 0, len(list))
	for _, item := range list {
		if item.ID == o.ID {
			continue
		}
		newList = append(newList, item)
	}
	roomOnlineCacheMgr.FlushCache(o.RoomId, newList)
}
