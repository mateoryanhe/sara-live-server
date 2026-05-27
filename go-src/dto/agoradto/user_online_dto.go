package agoradto

import "github.com/gogf/gf/v2/frame/g"

// CheckUserOnlineReq 根据房间ID与用户ID查询用户是否在声网频道内
type CheckUserOnlineReq struct {
	g.Meta `path:"/checkUserOnline" method:"post" summary:"查询用户是否在直播间(声网)" tags:"声网"`
	RoomId uint64 `json:"roomId" v:"required#房间ID不能为空" dc:"直播间ID"`
	UserId uint64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
}

// CheckUserOnlineRes 用户在线状态
type CheckUserOnlineRes struct {
	Online      bool   `json:"online" dc:"是否在频道内"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	UserAccount string `json:"userAccount" dc:"声网用户账号"`
	Role        int    `json:"role" dc:"用户在频道内角色(0未知,1通信用户,2主播,3观众),不在线时为0"`
	JoinAt      int64  `json:"joinAt" dc:"加入频道Unix秒,不在线时为0"`
}
