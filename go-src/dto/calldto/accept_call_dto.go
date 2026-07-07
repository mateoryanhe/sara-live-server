package calldto

import "github.com/gogf/gf/v2/frame/g"

// AcceptCallReq 同意接听通话
type AcceptCallReq struct {
	g.Meta  `path:"/acceptCall" method:"post" summary:"同意接听通话" tags:"通话"`
	OrderId uint64 `json:"orderId,string" v:"required|min:1#订单ID不能为空|订单ID无效" dc:"通话订单ID"`
}

// AcceptCallRes 同意接听通话响应
type AcceptCallRes struct {
	Success     bool   `json:"success"`
	OrderId     string `json:"orderId" dc:"通话订单ID"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	Token       string `json:"token" dc:"声网RTC Token"`
	AppId       string `json:"appId" dc:"声网AppId"`
	UserAccount string `json:"userAccount" dc:"声网用户账号(当前登录用户ID)"`
	ExpireAt    int64  `json:"expireAt" dc:"Token过期时间(Unix秒)"`
}
