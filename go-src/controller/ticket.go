package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/ticketdto"
	"xr-game-server/module/ticket"
)

const TicketUrl = "/ticket"

type TicketController struct{}

func initTicketController() {
	httpserver.RegCMS(TicketUrl, &TicketController{})
}

func (c *TicketController) TicketList(ctx context.Context, req *ticketdto.TicketListReq) (*httpserver.CMSQueryResp, error) {
	return ticket.GetTicketList(ctx, req)
}

func (c *TicketController) CreateTicket(ctx context.Context, req *ticketdto.CreateTicketReq) (*ticketdto.CreateTicketRes, error) {
	return ticket.CreateTicket(ctx, req)
}

func (c *TicketController) UpdateTicket(ctx context.Context, req *ticketdto.UpdateTicketReq) (*ticketdto.UpdateTicketRes, error) {
	return ticket.UpdateTicket(ctx, req)
}

func (c *TicketController) DeleteTicket(ctx context.Context, req *ticketdto.DeleteTicketReq) (*ticketdto.DeleteTicketRes, error) {
	return ticket.DeleteTicket(ctx, req)
}

func (c *TicketController) OnShelfTicket(ctx context.Context, req *ticketdto.OnShelfTicketReq) (*ticketdto.OnShelfTicketRes, error) {
	return ticket.OnShelfTicket(ctx, req)
}

func (c *TicketController) OffShelfTicket(ctx context.Context, req *ticketdto.OffShelfTicketReq) (*ticketdto.OffShelfTicketRes, error) {
	return ticket.OffShelfTicket(ctx, req)
}
