package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfTicketReq struct {
	g.Meta `path:"/onShelfTicket" method:"post" summary:"上架门票" tags:"门票"`
	ID     uint64 `json:"id" v:"required#门票ID不能为空" dc:"门票ID"`
}

type OnShelfTicketRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfTicketReq struct {
	g.Meta `path:"/offShelfTicket" method:"post" summary:"下架门票" tags:"门票"`
	ID     uint64 `json:"id" v:"required#门票ID不能为空" dc:"门票ID"`
}

type OffShelfTicketRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
