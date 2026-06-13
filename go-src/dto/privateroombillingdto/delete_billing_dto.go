package privateroombillingdto

import "github.com/gogf/gf/v2/frame/g"

type DeleteBillingReq struct {
	g.Meta `path:"/deleteBilling" method:"post" summary:"删除私密直播间计费" tags:"私密直播间计费"`
	ID     uint64 `json:"id" v:"required#计费配置ID不能为空" dc:"计费配置ID"`
}

type DeleteBillingRes struct {
	Success bool `json:"success"`
}
