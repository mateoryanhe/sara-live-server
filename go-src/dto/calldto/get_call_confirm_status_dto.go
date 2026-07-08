package calldto

import "github.com/gogf/gf/v2/frame/g"

// GetCallConfirmStatusReq 查询对方应答确认状态
type GetCallConfirmStatusReq struct {
	g.Meta  `path:"/getCallConfirmStatus" method:"post" summary:"查询对方应答确认状态" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// GetCallConfirmStatusRes 查询对方应答确认状态响应
type GetCallConfirmStatusRes struct {
	Success            bool   `json:"success"`
	OrderId            string `json:"orderId" dc:"通话订单ID"`
	SelfConfirmed      bool   `json:"selfConfirmed" dc:"当前用户是否已确认"`
	PeerConfirmed      bool   `json:"peerConfirmed" dc:"对方是否已确认"`
	PeerConfirmAt      int64  `json:"peerConfirmAt" dc:"对方确认时间(Unix秒,未确认时为0)"`
	AnswerConfirmCount uint32 `json:"answerConfirmCount" dc:"当前应答确认次数"`
	CallStarted        bool   `json:"callStarted" dc:"双方是否均已确认(通话已开始)"`
}
