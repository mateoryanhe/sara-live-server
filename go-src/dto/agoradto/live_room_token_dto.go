package agoradto

import "github.com/gogf/gf/v2/frame/g"

// GetLiveRoomTokenReq App端获取进直播间声网Token
type GetLiveRoomTokenReq struct {
	g.Meta `path:"/liveRoomToken" method:"post" summary:"获取进直播间声网Token" tags:"声网"`
	RoomId uint64 `json:"roomId" v:"required#房间ID不能为空" dc:"直播间ID"`
}

// GetLiveRoomTokenRes App端进直播间声网Token
type GetLiveRoomTokenRes struct {
	Token       string `json:"token" dc:"声网RTC Token"`
	AppId       string `json:"appId" dc:"声网AppId"`
	ChannelName string `json:"channelName" dc:"声网频道名(与roomId一致)"`
	UserAccount string `json:"userAccount" dc:"声网用户账号(与当前登录用户ID一致)"`
	ExpireAt    int64  `json:"expireAt" dc:"Token过期时间(Unix秒)"`
}
