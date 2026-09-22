package ticket

import (
	"context"

	"xr-game-server/dto/ticketdto"
)

func GetAppTicketList(_ context.Context, _ *ticketdto.AppTicketListReq) (*ticketdto.AppTicketListRes, error) {
	return &ticketdto.AppTicketListRes{List: getAppTicketList()}, nil
}
