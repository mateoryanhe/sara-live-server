package anchorvideopricetier

import (
	"sort"
	"strconv"
	"sync/atomic"

	"xr-game-server/dao/anchorvideopricetierdao"
	entity "xr-game-server/entity/live"
)

type PriceTierItem struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

type priceTierSnapshot struct {
	byID    map[uint64]*entity.AnchorVideoPriceTier
	onShelf []*PriceTierItem
}

var (
	priceTierCache     atomic.Value // *priceTierSnapshot
	emptyPriceTierList = make([]*PriceTierItem, 0)
)

func Init() {
	reloadMemory()
}

func reloadMemory() {
	rows := anchorvideopricetierdao.GetAll()
	byID := make(map[uint64]*entity.AnchorVideoPriceTier, len(rows))
	onShelf := make([]*PriceTierItem, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		byID[r.ID] = r
		if r.Status == entity.AnchorVideoPriceTierStatusOnShelf {
			onShelf = append(onShelf, &PriceTierItem{
				ID:    strconv.FormatUint(r.ID, 10),
				Price: r.Price,
			})
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].Price != onShelf[j].Price {
			return onShelf[i].Price < onShelf[j].Price
		}
		return onShelf[i].ID > onShelf[j].ID
	})
	priceTierCache.Store(&priceTierSnapshot{
		byID:    byID,
		onShelf: onShelf,
	})
}

func getSnapshot() *priceTierSnapshot {
	v := priceTierCache.Load()
	if v == nil {
		return &priceTierSnapshot{
			byID:    make(map[uint64]*entity.AnchorVideoPriceTier),
			onShelf: emptyPriceTierList,
		}
	}
	return v.(*priceTierSnapshot)
}

// GetOnShelfPriceTiers 上架挡位列表(仅内存;App 禁止查库).
func GetOnShelfPriceTiers() []*PriceTierItem {
	return getSnapshot().onShelf
}

// IsOnShelfPrice 价格是否命中已上架挡位(仅内存比较).
func IsOnShelfPrice(price float64) bool {
	for _, item := range GetOnShelfPriceTiers() {
		if item != nil && item.Price == price {
			return true
		}
	}
	return false
}
