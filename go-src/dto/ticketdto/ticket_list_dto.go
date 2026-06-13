package ticketdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type TicketListReq struct {
	g.Meta `path:"/ticketList" method:"post" summary:"获取门票列表" tags:"门票"`
	httpserver.CMSQueryReq
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=只看下架,2=只看上架)"`
}

type TicketListRes struct {
	ID        string  `json:"id"`
	Price     float64 `json:"price"`
	Sort      int     `json:"sort"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
