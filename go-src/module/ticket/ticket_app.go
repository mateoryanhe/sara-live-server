package ticket

import (
	"context"

	"xr-game-server/dto/ticketdto"
)

var emptyAppTicketList = make([]*ticketdto.AppTicketItem, 0)

// GetAppTicketList 已废弃：兼容旧 App，固定返回空列表，不再读 live_tickets。
func GetAppTicketList(_ context.Context, _ *ticketdto.AppTicketListReq) (*ticketdto.AppTicketListRes, error) {
	return &ticketdto.AppTicketListRes{List: emptyAppTicketList}, nil
}
