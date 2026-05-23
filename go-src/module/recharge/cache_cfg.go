package recharge

import (
	"sync/atomic"
	"xr-game-server/dao/rechargecfgdao"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity"
)

type rechargeCfgSnapshot struct {
	byId map[uint64]*rechargecfgdto.AppRechargeCfgItem
	list []*rechargecfgdto.AppRechargeCfgItem
}

// 充值配置缓存(仅缓存已上架配置,供 App 查询使用)
// 以 ID 作为主键,任何创建/修改/删除/上下架变更都会重新从 DB 加载并整体替换快照。
var (
	rechargeCfgCache     atomic.Value // *rechargeCfgSnapshot
	emptyRechargeCfgList = make([]*rechargecfgdto.AppRechargeCfgItem, 0)
)

// toAppItem 将 entity 转换为 App 端 DTO
func toAppItem(r *entity.RechargeCfg) *rechargecfgdto.AppRechargeCfgItem {
	return &rechargecfgdto.AppRechargeCfgItem{
		ID:          r.ID,
		Name:        r.Name,
		CfgType:     r.CfgType,
		Icon:        r.Icon,
		Gold:        r.Gold,
		ExtraGold:   r.ExtraGold,
		Price:       r.Price,
		Currency:    r.Currency,
		ProductId:   r.ProductId,
		Sort:        r.Sort,
		Description: r.Description,
		Status:      r.Status,
	}
}

// loadRechargeCfgCache 从 DB 重新加载已上架配置并整体替换快照
func loadRechargeCfgCache() []*rechargecfgdto.AppRechargeCfgItem {
	rows := rechargecfgdao.GetOnShelf()
	byId := make(map[uint64]*rechargecfgdto.AppRechargeCfgItem, len(rows))
	list := make([]*rechargecfgdto.AppRechargeCfgItem, 0, len(rows))
	for _, r := range rows {
		item := toAppItem(r)
		byId[r.ID] = item
		list = append(list, item)
	}
	rechargeCfgCache.Store(&rechargeCfgSnapshot{
		byId: byId,
		list: list,
	})
	return list
}

// reloadRechargeCfgCache CMS 端发生变更后重新加载缓存
func reloadRechargeCfgCache() {
	loadRechargeCfgCache()
}

func getRechargeCfgSnapshot() *rechargeCfgSnapshot {
	v := rechargeCfgCache.Load()
	if v == nil {
		return &rechargeCfgSnapshot{
			byId: make(map[uint64]*rechargecfgdto.AppRechargeCfgItem),
			list: emptyRechargeCfgList,
		}
	}
	return v.(*rechargeCfgSnapshot)
}

// getRechargeCfgCache 获取缓存列表,未加载则触发加载
func getRechargeCfgCache() []*rechargecfgdto.AppRechargeCfgItem {
	if rechargeCfgCache.Load() == nil {
		return loadRechargeCfgCache()
	}
	return getRechargeCfgSnapshot().list
}

// GetRechargeCfgFromCacheById 通过 ID 从缓存获取已上架的充值配置(供其它模块使用)
// 缓存只包含已上架配置;未加载则触发加载;命中失败返回 nil
func GetRechargeCfgFromCacheById(id uint64) *rechargecfgdto.AppRechargeCfgItem {
	if rechargeCfgCache.Load() == nil {
		loadRechargeCfgCache()
	}
	return getRechargeCfgSnapshot().byId[id]
}
