package ticketdto

import "github.com/gogf/gf/v2/frame/g"

type AppTicketListReq struct {
	g.Meta `path:"/appTicketList" method:"post" summary:"App查询门票列表(已上架)" tags:"门票"`
}

type AppTicketItem struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Sort  int     `json:"sort"`
}

type AppTicketListRes struct {
	List []*AppTicketItem `json:"list"`
}
