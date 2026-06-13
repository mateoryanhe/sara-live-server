package liveroomdao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// roomCacheMgr 按 roomId(== 主播用户ID) 缓存
var roomCacheMgr = gmap.NewKVMap[uint64, *entity.LiveRoom](false)

// InitLiveRoomDao 初始化直播间相关缓存
func initLiveRoomDao() {
	initLiveRoomTicketPayDao()
	//启动的时候,加载全部主播
	all := make([]*entity.LiveRoom, 0)
	g.Model(string(entity.TbLiveRoom)).Scan(&all)
	for _, v := range all {
		roomCacheMgr.Set(v.ID, v)
	}
}

// GetRoomById 按 roomId 获取直播间(走缓存)
func GetRoomById(roomId uint64) *entity.LiveRoom {
	if !roomCacheMgr.Contains(roomId) {
		return nil
	}
	v := roomCacheMgr.Get(roomId)

	return v
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

func GetAllLiveRoom() []*entity.LiveRoom {
	data := roomCacheMgr.Values()
	return data
}

// AddRoomToCache 新建直播间后写入缓存
func AddRoomToCache(r *entity.LiveRoom) {
	if r == nil {
		return
	}
	roomCacheMgr.Set(r.ID, r)
}

// FlushRoomCache 直播间字段变更后刷新缓存
func FlushRoomCache(r *entity.LiveRoom) {
	if r == nil {
		return
	}
	roomCacheMgr.Set(r.ID, r)
}
