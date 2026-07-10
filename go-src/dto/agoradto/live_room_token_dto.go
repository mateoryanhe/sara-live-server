package agoradto

import "github.com/gogf/gf/v2/frame/g"

const (
	RTCRolePublisher  uint8 = 1 // 主播/通话方(可发布音视频)
	RTCRoleSubscriber uint8 = 2 // 观众(仅订阅)
)

// GetLiveRoomTokenReq App端上报频道名获取声网Token(直播间或通话频道)
type GetLiveRoomTokenReq struct {
	g.Meta      `path:"/liveRoomToken" method:"post" summary:"获取声网频道Token" tags:"声网"`
	ChannelName string `json:"channelName" v:"required#频道名称不能为空" dc:"声网频道名"`
	Role        uint8  `json:"role" v:"required|in:1,2#RTC角色不能为空|RTC角色无效" dc:"RTC角色:1=Publisher 2=Subscriber"`
}

// GetLiveRoomTokenRes App端声网频道Token
type GetLiveRoomTokenRes struct {
	Token       string `json:"token" dc:"声网RTC Token"`
	AppId       string `json:"appId" dc:"声网AppId"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	UserAccount string `json:"userAccount" dc:"声网用户账号(与当前登录用户ID一致)"`
	ExpireAt    int64  `json:"expireAt" dc:"Token过期时间(Unix秒)"`
}
