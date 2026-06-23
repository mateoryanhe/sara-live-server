package shortvideo

import (
	"sort"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
)

type categorySnapshot struct {
	byID map[uint64]*entity.ShortVideoCategory
	list []*shortvideodto.AppShortVideoCategoryItem
}

var (
	categoryCache     atomic.Value // *categorySnapshot
	emptyCategoryList = make([]*shortvideodto.AppShortVideoCategoryItem, 0)
)

func initCategoryMemory() {
	reloadCategoryMemory()
}

func reloadCategoryMemory() {
	rows := shortvideodao.GetAllCategories()
	byID := make(map[uint64]*entity.ShortVideoCategory, len(rows))
	list := make([]*shortvideodto.AppShortVideoCategoryItem, 0, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
		list = append(list, toAppShortVideoCategoryItem(r))
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Sort != list[j].Sort {
			return list[i].Sort > list[j].Sort
		}
		return list[i].ID > list[j].ID
	})
	categoryCache.Store(&categorySnapshot{
		byID: byID,
		list: list,
	})
}

func getCategorySnapshot() *categorySnapshot {
	v := categoryCache.Load()
	if v == nil {
		return &categorySnapshot{
			byID: make(map[uint64]*entity.ShortVideoCategory),
			list: emptyCategoryList,
		}
	}
	return v.(*categorySnapshot)
}

func getAppCategoryList() []*shortvideodto.AppShortVideoCategoryItem {
	return getCategorySnapshot().list
}

func GetCategoryFromMemoryById(id uint64) *entity.ShortVideoCategory {
	return getCategorySnapshot().byID[id]
}

func toAppShortVideoCategoryItem(c *entity.ShortVideoCategory) *shortvideodto.AppShortVideoCategoryItem {
	return &shortvideodto.AppShortVideoCategoryItem{
		ID:   strconv.FormatUint(c.ID, 10),
		Name: c.Name,
		Sort: c.Sort,
	}
}
