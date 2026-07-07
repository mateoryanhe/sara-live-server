package calldto

import "github.com/gogf/gf/v2/frame/g"

// EndCallReq 结束通话
type EndCallReq struct {
	g.Meta  `path:"/endCall" method:"post" summary:"结束通话" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// EndCallRes 结束通话响应
type EndCallRes struct {
	Success         bool    `json:"success"`
	OrderId         string  `json:"orderId" dc:"通话订单ID"`
	CallDuration    uint32  `json:"callDuration" dc:"通话时长(秒)"`
	BillingDuration uint32  `json:"billingDuration" dc:"计费时长(分钟)"`
	TotalCost       float64 `json:"totalCost" dc:"总费用(钻石)"`
}
