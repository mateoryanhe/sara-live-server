package banner

import (
	"sort"
	"strconv"
	"sync"
	"xr-game-server/dao/bannerdao"
	"xr-game-server/dto/bannerdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

var (
	bannerMu            sync.RWMutex
	bannerMap           map[uint64]*entity.HomeBanner
	bannerOnShelfSorted []*bannerdto.AppBannerItem
	bannerLoaded        bool
)

func Init() {
	loadBannerMemory()
}

func loadBannerMemory() {
	rows := bannerdao.GetAll()
	m := make(map[uint64]*entity.HomeBanner, len(rows))
	onShelf := make([]*bannerdto.AppBannerItem, 0, len(rows))
	for _, r := range rows {
		m[r.ID] = r
		if r.Status == entity.HomeBannerStatusOnShelf {
			onShelf = append(onShelf, toAppBannerItem(r))
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].Sort != onShelf[j].Sort {
			return onShelf[i].Sort > onShelf[j].Sort
		}
		return onShelf[i].ID > onShelf[j].ID
	})

	bannerMu.Lock()
	bannerMap = m
	bannerOnShelfSorted = onShelf
	bannerLoaded = true
	bannerMu.Unlock()
}

func reloadBannerMemory() {
	loadBannerMemory()
}

func getAppBannerList() []*bannerdto.AppBannerItem {
	bannerMu.RLock()
	if bannerLoaded {
		list := bannerOnShelfSorted
		bannerMu.RUnlock()
		return list
	}
	bannerMu.RUnlock()
	loadBannerMemory()
	bannerMu.RLock()
	list := bannerOnShelfSorted
	bannerMu.RUnlock()
	return list
}

func GetBannerFromMemoryById(id uint64) *entity.HomeBanner {
	bannerMu.RLock()
	if bannerLoaded {
		row := bannerMap[id]
		bannerMu.RUnlock()
		return row
	}
	bannerMu.RUnlock()
	loadBannerMemory()
	bannerMu.RLock()
	row := bannerMap[id]
	bannerMu.RUnlock()
	return row
}

func toAppBannerItem(b *entity.HomeBanner) *bannerdto.AppBannerItem {
	return &bannerdto.AppBannerItem{
		ID:    strconv.FormatUint(b.ID, 10),
		Title: b.Title,
		Image: upload.GetUrlByName(b.Image),
		Link:  b.Link,
		Sort:  b.Sort,
	}
}
