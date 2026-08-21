package shortvideo

import (
	"sort"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
)

type priceTierSnapshot struct {
	byID    map[uint64]*entity.ShortVideoPriceTier
	onShelf []*shortvideodto.AppShortVideoPriceTierItem
}

var (
	priceTierCache     atomic.Value // *priceTierSnapshot
	emptyPriceTierList = make([]*shortvideodto.AppShortVideoPriceTierItem, 0)
)

func initPriceTierMemory() {
	reloadPriceTierMemory()
}

func reloadPriceTierMemory() {
	rows := shortvideodao.GetAllPriceTiers()
	byID := make(map[uint64]*entity.ShortVideoPriceTier, len(rows))
	onShelf := make([]*shortvideodto.AppShortVideoPriceTierItem, 0, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
		if r.Status == entity.ShortVideoPriceTierStatusOnShelf {
			onShelf = append(onShelf, toAppShortVideoPriceTierItem(r))
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

func getPriceTierSnapshot() *priceTierSnapshot {
	v := priceTierCache.Load()
	if v == nil {
		return &priceTierSnapshot{
			byID:    make(map[uint64]*entity.ShortVideoPriceTier),
			onShelf: emptyPriceTierList,
		}
	}
	return v.(*priceTierSnapshot)
}

func getAppPriceTierList() []*shortvideodto.AppShortVideoPriceTierItem {
	return getPriceTierSnapshot().onShelf
}

func toAppShortVideoPriceTierItem(t *entity.ShortVideoPriceTier) *shortvideodto.AppShortVideoPriceTierItem {
	return &shortvideodto.AppShortVideoPriceTierItem{
		ID:    strconv.FormatUint(t.ID, 10),
		Price: t.Price,
	}
}
