package calldto

import "github.com/gogf/gf/v2/frame/g"

// LiveRoomCallReq 直播间通话呼叫
type LiveRoomCallReq struct {
	g.Meta   `path:"/liveRoomCall" method:"post" summary:"直播间通话呼叫" tags:"通话"`
	AnchorId uint64 `json:"anchorId" v:"required|min:1#主播ID不能为空|主播ID无效" dc:"主播ID(即直播间ID)"`
}

// LiveRoomCallRes 直播间通话呼叫响应
type LiveRoomCallRes struct {
	OrderId     string `json:"orderId" dc:"通话订单ID"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	Token       string `json:"token" dc:"声网RTC Token"`
	AppId       string `json:"appId" dc:"声网AppId"`
	UserAccount string `json:"userAccount" dc:"声网用户账号(当前登录用户ID)"`
	ExpireAt    int64  `json:"expireAt" dc:"Token过期时间(Unix秒)"`
}
