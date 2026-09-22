package ticket

import (
	"sort"
	"strconv"
	"sync/atomic"

	"xr-game-server/dao/ticketdao"
	"xr-game-server/dto/ticketdto"
	liveentity "xr-game-server/entity/live"
)

type ticketSnapshot struct {
	onShelf []*ticketdto.AppTicketItem
}

var (
	ticketCache     atomic.Value // *ticketSnapshot
	emptyTicketList = make([]*ticketdto.AppTicketItem, 0)
)

func Init() {
	reloadTicketMemory()
}

// reloadTicketMemory 从数据库加载已上架门票并整体替换内存快照。
func reloadTicketMemory() {
	rows := ticketdao.GetAll()
	onShelf := make([]*ticketdto.AppTicketItem, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.Status == liveentity.LiveTicketStatusOnShelf {
			onShelf = append(onShelf, toAppTicketItem(row))
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].Sort != onShelf[j].Sort {
			return onShelf[i].Sort > onShelf[j].Sort
		}
		return onShelf[i].ID > onShelf[j].ID
	})
	ticketCache.Store(&ticketSnapshot{onShelf: onShelf})
}

func getAppTicketList() []*ticketdto.AppTicketItem {
	v := ticketCache.Load()
	if v == nil {
		return emptyTicketList
	}
	return v.(*ticketSnapshot).onShelf
}

func toAppTicketItem(ticket *liveentity.LiveTicket) *ticketdto.AppTicketItem {
	return &ticketdto.AppTicketItem{
		ID:    strconv.FormatUint(ticket.ID, 10),
		Price: ticket.Price,
		Sort:  ticket.Sort,
	}
}
