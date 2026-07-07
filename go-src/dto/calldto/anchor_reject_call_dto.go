package calldto

import "github.com/gogf/gf/v2/frame/g"

// AnchorRejectCallReq 主播拒接通话
type AnchorRejectCallReq struct {
	g.Meta  `path:"/anchorRejectCall" method:"post" summary:"拒接通话" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// AnchorRejectCallRes 主播拒接通话响应
type AnchorRejectCallRes struct {
	Success bool `json:"success"`
}
