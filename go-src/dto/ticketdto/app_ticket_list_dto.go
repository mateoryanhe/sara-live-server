package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type AppTicketListReq struct {
	g.Meta `path:"/appTicketList" method:"post" summary:"[已废弃]固定返回空列表,不再读库" tags:"门票"`
}

type AppTicketItem struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Sort  int     `json:"sort"`
}

type AppTicketListRes struct {
	List []*AppTicketItem `json:"list"`
}
