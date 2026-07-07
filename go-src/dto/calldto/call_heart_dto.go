package calldto

import "github.com/gogf/gf/v2/frame/g"

// CallHeartReq 通话心跳(呼叫者与接听者每10秒调用一次)
type CallHeartReq struct {
	g.Meta  `path:"/callHeart" method:"post" summary:"通话心跳" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// CallHeartRes 通话心跳响应
type CallHeartRes struct {
	Success bool `json:"success"`
}
