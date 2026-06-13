package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type CreateTicketReq struct {
	g.Meta `path:"/createTicket" method:"post" summary:"创建门票" tags:"门票"`
	Price  float64 `json:"price" dc:"钻石价格"`
	Sort   int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateTicketRes struct {
	ID string `json:"id" dc:"门票ID"`
}
