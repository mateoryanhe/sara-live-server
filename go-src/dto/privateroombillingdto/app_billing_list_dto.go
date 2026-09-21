package privateroombillingdto

import "github.com/gogf/gf/v2/frame/g"

type AppBillingListReq struct {
	g.Meta `path:"/appBillingList" method:"post" summary:"App查询1v1视频通话计费列表(已上架)" tags:"1v1视频通话计费"`
}

type AppBillingItem struct {
	ID             string  `json:"id"`
	PricePerMinute float64 `json:"pricePerMinute"`
	Sort           int     `json:"sort"`
}

type AppBillingListRes struct {
	List []*AppBillingItem `json:"list"`
}
