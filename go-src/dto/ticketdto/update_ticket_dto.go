package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type UpdateTicketReq struct {
	g.Meta `path:"/updateTicket" method:"post" summary:"修改门票" tags:"门票"`
	ID     uint64  `json:"id" v:"required#门票ID不能为空" dc:"门票ID"`
	Price  float64 `json:"price" dc:"钻石价格"`
	Sort   int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateTicketRes struct {
	Success bool `json:"success"`
}
