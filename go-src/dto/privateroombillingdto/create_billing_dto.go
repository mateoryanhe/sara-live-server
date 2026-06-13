package privateroombillingdto

import "github.com/gogf/gf/v2/frame/g"

type CreateBillingReq struct {
	g.Meta         `path:"/createBilling" method:"post" summary:"创建私密直播间计费" tags:"私密直播间计费"`
	PricePerMinute float64 `json:"pricePerMinute" dc:"每分钟钻石价格"`
	Sort           int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateBillingRes struct {
	ID string `json:"id" dc:"计费配置ID"`
}
