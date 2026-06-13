package privateroombillingdto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfBillingReq struct {
	g.Meta `path:"/onShelfBilling" method:"post" summary:"上架私密直播间计费" tags:"私密直播间计费"`
	ID     uint64 `json:"id" v:"required#计费配置ID不能为空" dc:"计费配置ID"`
}

type OnShelfBillingRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfBillingReq struct {
	g.Meta `path:"/offShelfBilling" method:"post" summary:"下架私密直播间计费" tags:"私密直播间计费"`
	ID     uint64 `json:"id" v:"required#计费配置ID不能为空" dc:"计费配置ID"`
}

type OffShelfBillingRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
