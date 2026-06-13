package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type DeleteTicketReq struct {
	g.Meta `path:"/deleteTicket" method:"post" summary:"删除门票" tags:"门票"`
	ID     uint64 `json:"id" v:"required#门票ID不能为空" dc:"门票ID"`
}

type DeleteTicketRes struct {
	Success bool `json:"success"`
}
