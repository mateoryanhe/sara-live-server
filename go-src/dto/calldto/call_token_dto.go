package calldto

import "github.com/gogf/gf/v2/frame/g"

// CallTokenReq App端获取通话频道Token
type CallTokenReq struct {
	g.Meta      `path:"/callToken" method:"post" summary:"获取通话频道Token" tags:"通话"`
	ChannelName string `json:"channelName" v:"required#频道名称不能为空" dc:"声网频道名"`
}

// CallTokenRes App端通话频道Token
type CallTokenRes struct {
	Token       string `json:"token" dc:"声网RTC Token"`
	AppId       string `json:"appId" dc:"声网AppId"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	UserAccount string `json:"userAccount" dc:"声网用户账号(当前登录用户ID)"`
	ExpireAt    int64  `json:"expireAt" dc:"Token过期时间(Unix秒)"`
}
