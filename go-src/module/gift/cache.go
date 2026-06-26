package gift

import (
	"sync"
	"xr-game-server/dao/giftdao"
	"xr-game-server/dto/giftdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

// 礼物缓存(仅缓存已上架礼物,供 App 查询使用)
// 以礼物 ID 作为主键,任何创建/修改/删除/上下架变更都会清空缓存,下一次读取时重新加载。
var (
	giftCacheMu     sync.RWMutex
	giftCacheMap    map[uint64]*giftdto.AppGiftItem // id -> item
	giftCacheSorted []*giftdto.AppGiftItem          // 按 sort desc, created_at desc 排序的列表
	giftCacheLoaded bool
)

// invalidateGiftCache 清除礼物缓存
func invalidateGiftCache() {
	giftCacheMu.Lock()
	giftCacheMap = nil
	giftCacheSorted = nil
	giftCacheLoaded = false
	giftCacheMu.Unlock()
}

// loadGiftCache 从 DB 重新加载已上架礼物到缓存
func loadGiftCache() []*giftdto.AppGiftItem {
	rows := giftdao.GetOnShelfGifts()
	m := make(map[uint64]*giftdto.AppGiftItem, len(rows))
	list := make([]*giftdto.AppGiftItem, 0, len(rows))
	for _, r := range rows {
		item := toAppGiftItem(r)
		m[r.ID] = item
		list = append(list, item)
	}

	giftCacheMu.Lock()
	giftCacheMap = m
	giftCacheSorted = list
	giftCacheLoaded = true
	giftCacheMu.Unlock()
	return list
}

// getGiftCache 获取缓存列表,未加载则触发加载
func getGiftCache() []*giftdto.AppGiftItem {
	giftCacheMu.RLock()
	if giftCacheLoaded {
		list := giftCacheSorted
		giftCacheMu.RUnlock()
		return list
	}
	giftCacheMu.RUnlock()
	return loadGiftCache()
}

// GetGiftFromCacheById 通过礼物 ID 从缓存获取礼物配置(供其它模块使用,例如送礼)
// 缓存只包含已上架礼物;若未加载则触发加载;命中失败返回 nil
func GetGiftFromCacheById(id uint64) *giftdto.AppGiftItem {
	giftCacheMu.RLock()
	if giftCacheLoaded {
		item := giftCacheMap[id]
		giftCacheMu.RUnlock()
		return item
	}
	giftCacheMu.RUnlock()
	loadGiftCache()
	giftCacheMu.RLock()
	item := giftCacheMap[id]
	giftCacheMu.RUnlock()
	return item
}

func toAppGiftItem(g *entity.LiveGift) *giftdto.AppGiftItem {
	var publishedAt int64
	if g.PublishedAt != nil {
		publishedAt = g.PublishedAt.Unix()
	}
	return &giftdto.AppGiftItem{
		ID:          g.ID,
		Name:        g.Name,
		NameEn:      g.NameEn,
		NameEs:      g.NameEs,
		NamePt:      g.NamePt,
		NameHi:      g.NameHi,
		Icon:        upload.GetUrlByName(g.Icon),
		Animation:   upload.GetUrlByName(g.Animation),
		Price:       g.Price,
		Category:    g.Category,
		Sort:        g.Sort,
		PublishedAt: publishedAt,
		Description: g.Description,
	}
}
