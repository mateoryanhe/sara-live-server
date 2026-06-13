package privateroombillingdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type BillingListReq struct {
	g.Meta `path:"/billingList" method:"post" summary:"获取私密直播间计费列表" tags:"私密直播间计费"`
	httpserver.CMSQueryReq
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=只看下架,2=只看上架)"`
}

type BillingListRes struct {
	ID             string  `json:"id"`
	PricePerMinute float64 `json:"pricePerMinute"`
	Sort           int     `json:"sort"`
	Status         uint8   `json:"status"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}
