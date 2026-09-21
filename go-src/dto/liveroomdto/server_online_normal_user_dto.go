package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// GetServerOnlineNormalUserListReq 查询当前服务器在线普通用户快照。
type GetServerOnlineNormalUserListReq struct {
	g.Meta `path:"/serverOnlineNormalUserList" method:"post" summary:"查询当前服务器在线普通用户" tags:"直播间"`
}

// ServerOnlineNormalUserItem 当前服务器在线普通用户基础资料。
type ServerOnlineNormalUserItem struct {
	UserId   string `json:"userId" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像URL(已拼资源域名)"`
	VipLevel uint32 `json:"vipLevel" dc:"VIP等级"`
	Gender   uint8  `json:"gender" dc:"性别(0未知,1男,2女)"`
	Age      int    `json:"age" dc:"年龄(未设置出生日期时为0)"`
	UserType uint8  `json:"userType" dc:"用户类型(固定为0普通用户)"`
}

type GetServerOnlineNormalUserListRes struct {
	Total       int                           `json:"total" dc:"在线普通用户数量(最大500)"`
	List        []*ServerOnlineNormalUserItem `json:"list" dc:"在线普通用户列表"`
	RefreshedAt int64                         `json:"refreshedAt" dc:"缓存刷新时间(毫秒)"`
	SysTime     int64                         `json:"sysTime" dc:"服务器当前时间(毫秒)"`
}
