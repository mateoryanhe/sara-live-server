package liveroomdao

import (
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

// roomCacheMgr 按 roomId(== 主播用户ID) 缓存
var roomCacheMgr = gmap.NewKVMap[uint64, *entity.LiveRoom](false)

// InitLiveRoomDao 初始化直播间相关缓存
func initLiveRoomDao() {
	initLiveRoomBillingPayDao()
	// 启动时只加载上架(status=1)的直播间
	all := make([]*entity.LiveRoom, 0)
	_ = g.Model(string(entity.TbLiveRoom)).
		Where(string(entity.LiveRoomStatus), entity.LiveRoomStatusOnShelf).
		Scan(&all)
	for _, v := range all {
		roomCacheMgr.Set(v.ID, v)
	}
	ids := make([]uint64, 0, len(all))
	for _, v := range all {
		if v != nil && v.ID != 0 {
			ids = append(ids, v.ID)
		}
	}
	PreloadLiveRoomIncomes(ids)
	PreloadLiveRoomCfgs(ids)
}

// GetRoomFromDB 按 roomId 直查数据库(含下架直播间,不走缓存)
func GetRoomFromDB(roomId uint64) *entity.LiveRoom {
	if roomId == 0 {
		return nil
	}
	var row entity.LiveRoom
	err := g.Model(string(entity.TbLiveRoom)).WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetRoomById 按 roomId 获取直播间(走缓存,仅上架)
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

// GetAnchorGuildId 主播所属工会ID(live_room.guild_id; roomId==用户ID; 无直播间返回0)
func GetAnchorGuildId(anchorId uint64) uint64 {
	if room := GetRoomById(anchorId); room != nil {
		return room.GuildId
	}
	if room := GetRoomFromDB(anchorId); room != nil {
		return room.GuildId
	}
	return 0
}

// ResolveRoom 优先读上架缓存,否则直查DB(含下架)
func ResolveRoom(anchorId uint64) *entity.LiveRoom {
	if room := GetRoomById(anchorId); room != nil {
		return room
	}
	return GetRoomFromDB(anchorId)
}

// HasLiveRoom 是否已有直播间记录(roomId==用户ID)
func HasLiveRoom(anchorId uint64) bool {
	return ResolveRoom(anchorId) != nil
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

// ListActiveCloudPlayerRooms 查询正在直播且已创建云播放器的直播间
func ListActiveCloudPlayerRooms() []*entity.LiveRoom {
	rooms := make([]*entity.LiveRoom, 0)
	for _, room := range roomCacheMgr.Values() {
		if room == nil || room.ID == 0 || room.LiveRecordId == 0 {
			continue
		}
		cfg := GetLiveRoomCfgFromCache(room.ID)
		if cfg == nil || strings.TrimSpace(cfg.CloudPlayerId) == "" {
			continue
		}
		rooms = append(rooms, room)
	}
	return rooms
}

func GetAllLiveRoom() []*entity.LiveRoom {
	data := roomCacheMgr.Values()
	return data
}

// FlushRoomCache 直播间字段变更后刷新缓存(仅上架房间保留在缓存)
func FlushRoomCache(r *entity.LiveRoom) {
	if r == nil {
		return
	}
	if r.Status != entity.LiveRoomStatusOnShelf {
		roomCacheMgr.Remove(r.ID)
		RemoveLiveRoomIncomeFromCache(r.ID)
		RemoveLiveRoomCfgFromCache(r.ID)
		return
	}
	roomCacheMgr.Set(r.ID, r)
}

// AddRoomToCache 新建/上架直播间后写入缓存(仅上架)
func AddRoomToCache(r *entity.LiveRoom) {
	if r == nil || r.Status != entity.LiveRoomStatusOnShelf {
		return
	}
	roomCacheMgr.Set(r.ID, r)
	AddLiveRoomIncomeToCache(r.ID)
	AddLiveRoomCfgToCache(r.ID)
}

// RemoveRoomFromCache 从直播间缓存移除(停用机器人主播等场景)
func RemoveRoomFromCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	roomCacheMgr.Remove(roomId)
	RemoveLiveRoomIncomeFromCache(roomId)
	RemoveLiveRoomCfgFromCache(roomId)
}

// ListRoomsByGuild 获取指定工会下的所有直播间(缓存,仅上架)
func ListRoomsByGuild(guildId uint64) []*entity.LiveRoom {
	rooms := make([]*entity.LiveRoom, 0)
	if guildId == 0 {
		return rooms
	}
	for _, room := range roomCacheMgr.Values() {
		if room == nil || room.ID == 0 {
			continue
		}
		if room.GuildId == guildId {
			rooms = append(rooms, room)
		}
	}
	return rooms
}

// ListRoomsByGuildFromDB 直查指定工会下全部直播间(含下架)
func ListRoomsByGuildFromDB(guildId uint64) []*entity.LiveRoom {
	rooms := make([]*entity.LiveRoom, 0)
	if guildId == 0 {
		return rooms
	}
	_ = g.Model(string(entity.TbLiveRoom)).
		Where(string(entity.LiveRoomGuildId), guildId).
		Scan(&rooms)
	return rooms
}

// GetGuildIdMapByRoomIds 批量查询直播间所属工会ID
func GetGuildIdMapByRoomIds(roomIds []uint64) map[uint64]uint64 {
	if len(roomIds) == 0 {
		return nil
	}
	type row struct {
		ID      uint64
		GuildId uint64
	}
	rows := make([]row, 0)
	_ = g.Model(string(entity.TbLiveRoom)).
		Where("id IN (?)", roomIds).
		Scan(&rows)
	ret := make(map[uint64]uint64, len(rows))
	for _, r := range rows {
		ret[r.ID] = r.GuildId
	}
	return ret
}
