package privateroombilling

import (
	"sort"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/privateroombillingdao"
	"xr-game-server/dto/privateroombillingdto"
	"xr-game-server/entity"
)

type billingSnapshot struct {
	byID    map[uint64]*entity.LivePrivateRoomBilling
	onShelf []*privateroombillingdto.AppBillingItem
}

var (
	billingCache     atomic.Value // *billingSnapshot
	emptyBillingList = make([]*privateroombillingdto.AppBillingItem, 0)
)

func Init() {
	reloadBillingMemory()
}

// reloadBillingMemory 从DB重新加载并整体替换内存快照
func reloadBillingMemory() {
	rows := privateroombillingdao.GetAll()
	byID := make(map[uint64]*entity.LivePrivateRoomBilling, len(rows))
	onShelf := make([]*privateroombillingdto.AppBillingItem, 0, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
		if r.Status == entity.LivePrivateRoomBillingStatusOnShelf {
			onShelf = append(onShelf, toAppBillingItem(r))
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].Sort != onShelf[j].Sort {
			return onShelf[i].Sort > onShelf[j].Sort
		}
		return onShelf[i].ID > onShelf[j].ID
	})
	billingCache.Store(&billingSnapshot{
		byID:    byID,
		onShelf: onShelf,
	})
}

func getBillingSnapshot() *billingSnapshot {
	v := billingCache.Load()
	if v == nil {
		return &billingSnapshot{
			byID:    make(map[uint64]*entity.LivePrivateRoomBilling),
			onShelf: emptyBillingList,
		}
	}
	return v.(*billingSnapshot)
}

func getAppBillingList() []*privateroombillingdto.AppBillingItem {
	return getBillingSnapshot().onShelf
}

func toAppBillingItem(row *entity.LivePrivateRoomBilling) *privateroombillingdto.AppBillingItem {
	return &privateroombillingdto.AppBillingItem{
		ID:             strconv.FormatUint(row.ID, 10),
		PricePerMinute: row.PricePerMinute,
		Sort:           row.Sort,
	}
}
