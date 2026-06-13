package ticket

import (
	"sort"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/ticketdao"
	"xr-game-server/dto/ticketdto"
	"xr-game-server/entity"
)

type ticketSnapshot struct {
	byID    map[uint64]*entity.LiveTicket
	onShelf []*ticketdto.AppTicketItem
}

var (
	ticketCache     atomic.Value // *ticketSnapshot
	emptyTicketList = make([]*ticketdto.AppTicketItem, 0)
)

func Init() {
	reloadTicketMemory()
}

// reloadTicketMemory 从DB重新加载并整体替换内存快照
func reloadTicketMemory() {
	rows := ticketdao.GetAll()
	byID := make(map[uint64]*entity.LiveTicket, len(rows))
	onShelf := make([]*ticketdto.AppTicketItem, 0, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
		if r.Status == entity.LiveTicketStatusOnShelf {
			onShelf = append(onShelf, toAppTicketItem(r))
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].Sort != onShelf[j].Sort {
			return onShelf[i].Sort > onShelf[j].Sort
		}
		return onShelf[i].ID > onShelf[j].ID
	})
	ticketCache.Store(&ticketSnapshot{
		byID:    byID,
		onShelf: onShelf,
	})
}

func getTicketSnapshot() *ticketSnapshot {
	v := ticketCache.Load()
	if v == nil {
		return &ticketSnapshot{
			byID:    make(map[uint64]*entity.LiveTicket),
			onShelf: emptyTicketList,
		}
	}
	return v.(*ticketSnapshot)
}

func getAppTicketList() []*ticketdto.AppTicketItem {
	return getTicketSnapshot().onShelf
}

func toAppTicketItem(t *entity.LiveTicket) *ticketdto.AppTicketItem {
	return &ticketdto.AppTicketItem{
		ID:    strconv.FormatUint(t.ID, 10),
		Price: t.Price,
		Sort:  t.Sort,
	}
}
