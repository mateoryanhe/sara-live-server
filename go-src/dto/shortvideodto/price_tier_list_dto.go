package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoPriceTierListReq struct {
	g.Meta `path:"/shortVideoPriceTierList" method:"post" summary:"获取短视频价格挡位列表" tags:"短视频"`
	httpserver.CMSQueryReq
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=只看下架,2=只看上架)"`
}

type ShortVideoPriceTierListRes struct {
	ID        string  `json:"id"`
	Price     float64 `json:"price"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
