package recharge

import (
	"strconv"
	"sync"
	"xr-game-server/dao/rechargecfgdao"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity"
)

// 充值配置缓存(仅缓存已上架配置,供 App 查询使用)
// 以 ID 作为主键,任何创建/修改/删除/上下架变更都会重新从 DB 加载缓存。
var (
	rechargeCfgCacheMu     sync.RWMutex
	rechargeCfgCacheMap    map[uint64]*rechargecfgdto.AppRechargeCfgItem // id -> item
	rechargeCfgCacheSorted []*rechargecfgdto.AppRechargeCfgItem          // 已排序列表
	rechargeCfgCacheLoaded bool
)

// Init 服务启动时预加载充值配置缓存
func Init() {
	loadRechargeCfgCache()
}

// toAppItem 将 entity 转换为 App 端 DTO
func toAppItem(r *entity.RechargeCfg) *rechargecfgdto.AppRechargeCfgItem {
	return &rechargecfgdto.AppRechargeCfgItem{
		ID:          strconv.FormatUint(r.ID, 10),
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
	}
}

// loadRechargeCfgCache 从 DB 重新加载已上架配置到缓存
func loadRechargeCfgCache() []*rechargecfgdto.AppRechargeCfgItem {
	rows := rechargecfgdao.GetOnShelf()
	m := make(map[uint64]*rechargecfgdto.AppRechargeCfgItem, len(rows))
	list := make([]*rechargecfgdto.AppRechargeCfgItem, 0, len(rows))
	for _, r := range rows {
		item := toAppItem(r)
		m[r.ID] = item
		list = append(list, item)
	}

	rechargeCfgCacheMu.Lock()
	rechargeCfgCacheMap = m
	rechargeCfgCacheSorted = list
	rechargeCfgCacheLoaded = true
	rechargeCfgCacheMu.Unlock()
	return list
}

// reloadRechargeCfgCache CMS 端发生变更后重新加载缓存
func reloadRechargeCfgCache() {
	loadRechargeCfgCache()
}

// getRechargeCfgCache 获取缓存列表,未加载则触发加载
func getRechargeCfgCache() []*rechargecfgdto.AppRechargeCfgItem {
	rechargeCfgCacheMu.RLock()
	if rechargeCfgCacheLoaded {
		list := rechargeCfgCacheSorted
		rechargeCfgCacheMu.RUnlock()
		return list
	}
	rechargeCfgCacheMu.RUnlock()
	return loadRechargeCfgCache()
}

// GetRechargeCfgFromCacheById 通过 ID 从缓存获取已上架的充值配置(供其它模块使用)
// 缓存只包含已上架配置;未加载则触发加载;命中失败返回 nil
func GetRechargeCfgFromCacheById(id uint64) *rechargecfgdto.AppRechargeCfgItem {
	rechargeCfgCacheMu.RLock()
	if rechargeCfgCacheLoaded {
		item := rechargeCfgCacheMap[id]
		rechargeCfgCacheMu.RUnlock()
		return item
	}
	rechargeCfgCacheMu.RUnlock()
	loadRechargeCfgCache()
	rechargeCfgCacheMu.RLock()
	item := rechargeCfgCacheMap[id]
	rechargeCfgCacheMu.RUnlock()
	return item
}
