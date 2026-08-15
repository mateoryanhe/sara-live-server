package liveroomdao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/live"
)

var liveRoomCfgCache = gmap.NewKVMap[uint64, *entity.LiveRoomCfg](false)

// PreloadLiveRoomCfgs 预热上架直播间配置(永久缓存)
func PreloadLiveRoomCfgs(roomIds []uint64) {
	if len(roomIds) == 0 {
		return
	}
	rows := make([]*entity.LiveRoomCfg, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomCfg)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&rows)
	rowMap := make(map[uint64]*entity.LiveRoomCfg, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			rowMap[row.ID] = row
		}
	}
	for _, id := range roomIds {
		if id == 0 {
			continue
		}
		if row, ok := rowMap[id]; ok {
			liveRoomCfgCache.Set(id, row)
		} else {
			liveRoomCfgCache.Set(id, entity.NewLiveRoomCfg(id))
		}
	}
}

// AddLiveRoomCfgToCache 上架时加载/创建配置缓存
func AddLiveRoomCfgToCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	PreloadLiveRoomCfgs([]uint64{roomId})
}

// RemoveLiveRoomCfgFromCache 下架时移除配置缓存
func RemoveLiveRoomCfgFromCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	liveRoomCfgCache.Remove(roomId)
}

// GetLiveRoomCfg 获取直播间配置(仅内存,无则新建并写入缓存)
func GetLiveRoomCfg(roomId uint64) *entity.LiveRoomCfg {
	if roomId == 0 {
		return nil
	}
	if liveRoomCfgCache.Contains(roomId) {
		return liveRoomCfgCache.Get(roomId)
	}
	row := entity.NewLiveRoomCfg(roomId)
	liveRoomCfgCache.Set(roomId, row)
	return row
}

// GetLiveRoomCfgFromCache 仅读缓存,未命中返回 nil
func GetLiveRoomCfgFromCache(roomId uint64) *entity.LiveRoomCfg {
	if roomId == 0 || !liveRoomCfgCache.Contains(roomId) {
		return nil
	}
	return liveRoomCfgCache.Get(roomId)
}

// GetLiveRoomCfgFromDB 直查数据库(回收站等场景)
func GetLiveRoomCfgFromDB(roomId uint64) *entity.LiveRoomCfg {
	if roomId == 0 {
		return nil
	}
	var row entity.LiveRoomCfg
	err := g.Model(string(entity.TbLiveRoomCfg)).Unscoped().WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListLiveRoomCfgFromDB 批量直查配置
func ListLiveRoomCfgFromDB(roomIds []uint64) map[uint64]*entity.LiveRoomCfg {
	ret := make(map[uint64]*entity.LiveRoomCfg, len(roomIds))
	if len(roomIds) == 0 {
		return ret
	}
	rows := make([]*entity.LiveRoomCfg, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomCfg)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			ret[row.ID] = row
		}
	}
	return ret
}
