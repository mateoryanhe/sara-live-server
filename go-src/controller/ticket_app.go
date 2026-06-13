package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/ticketdto"
	"xr-game-server/module/ticket"
)

const TicketAppUrl = "/ticket"

type TicketAppController struct{}

func initTicketAppController() {
	httpserver.RegAPI(TicketAppUrl, &TicketAppController{})
}

func (c *TicketAppController) AppTicketList(ctx context.Context, req *ticketdto.AppTicketListReq) (*ticketdto.AppTicketListRes, error) {
	return ticket.GetAppTicketList(ctx, req)
}
