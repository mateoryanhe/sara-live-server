package privateroombillingdto

import "github.com/gogf/gf/v2/frame/g"

type UpdateBillingReq struct {
	g.Meta         `path:"/updateBilling" method:"post" summary:"修改1v1视频通话计费" tags:"1v1视频通话计费"`
	ID             uint64  `json:"id" v:"required#计费配置ID不能为空" dc:"计费配置ID"`
	PricePerMinute float64 `json:"pricePerMinute" dc:"每分钟钻石价格"`
	Sort           int     `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateBillingRes struct {
	Success bool `json:"success"`
}
