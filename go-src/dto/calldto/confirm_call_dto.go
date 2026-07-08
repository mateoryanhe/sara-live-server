package calldto

import "github.com/gogf/gf/v2/frame/g"

// ConfirmCallReq 通话应答确认(双方各确认一次,累计2次后发起首次扣费)
type ConfirmCallReq struct {
	g.Meta  `path:"/confirmCall" method:"post" summary:"通话应答确认" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// ConfirmCallRes 通话应答确认响应
type ConfirmCallRes struct {
	Success             bool   `json:"success"`
	OrderId             string `json:"orderId" dc:"通话订单ID"`
	AnswerConfirmCount  uint32 `json:"answerConfirmCount" dc:"当前应答确认次数"`
	FirstChargeExecuted bool   `json:"firstChargeExecuted" dc:"本次是否已执行首次扣费"`
}
