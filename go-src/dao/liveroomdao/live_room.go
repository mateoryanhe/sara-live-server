package liveroomdao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

// roomCacheMgr 按 roomId(== 主播用户ID) 缓存
var roomCacheMgr *cache.CacheMgr

// InitLiveRoomDao 初始化直播间相关缓存
func InitLiveRoomDao() {
	roomCacheMgr = cache.NewCacheMgr()
}

// GetRoomById 按 roomId 获取直播间(走缓存)
func GetRoomById(roomId uint64) *entity.LiveRoom {
	v := roomCacheMgr.GetData(roomId, func(ctx context.Context) (value interface{}, err error) {
		var r *entity.LiveRoom
		_ = g.Model(string(entity.TbLiveRoom)).Where("id = ?", roomId).Scan(&r)
		return r, nil
	})
	if v == nil {
		return nil
	}
	r, _ := v.(*entity.LiveRoom)
	return r
}

// GetRoomByAnchor 按主播ID(== roomId)获取直播间;保留函数名以明语义
func GetRoomByAnchor(anchorId uint64) *entity.LiveRoom {
	return GetRoomById(anchorId)
}

// ListLivingRoomIds 查询正在直播的直播间ID(live_record_id > 0),仅返回 id 列
func ListLivingRoomIds() []uint64 {
	type idRow struct {
		ID uint64 `json:"id"`
	}
	rows := make([]*idRow, 0)
	_ = g.Model(string(entity.TbLiveRoom)).
		Fields("id").
		Where("live_record_id > ?", 0).
		Scan(&rows)
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ids = append(ids, row.ID)
	}
	return ids
}

// AddRoomToCache 新建直播间后写入缓存
func AddRoomToCache(r *entity.LiveRoom) {
	if r == nil {
		return
	}
	roomCacheMgr.FlushCache(r.ID, r)
}

// FlushRoomCache 直播间字段变更后刷新缓存
func FlushRoomCache(r *entity.LiveRoom) {
	if r == nil {
		return
	}
	roomCacheMgr.FlushCache(r.ID, r)
}

// RemoveRoomCache 移除指定直播间的缓存
func RemoveRoomCache(roomId uint64) {
	if roomCacheMgr == nil {
		return
	}
	_, _ = roomCacheMgr.Cache.Remove(gctx.New(), roomId)
}
