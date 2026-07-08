package calldto

import "github.com/gogf/gf/v2/frame/g"

// CallTimeoutReq 呼叫超时
type CallTimeoutReq struct {
	g.Meta  `path:"/callTimeout" method:"post" summary:"呼叫超时" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// CallTimeoutRes 呼叫超时响应
type CallTimeoutRes struct {
	Success bool `json:"success"`
}
