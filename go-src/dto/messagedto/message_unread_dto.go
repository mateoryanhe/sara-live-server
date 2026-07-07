package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppMessageUnreadCountReq App端查询消息未读数
type AppMessageUnreadCountReq struct {
	g.Meta `path:"/messageUnreadCount" method:"post" summary:"查询消息未读数" tags:"消息"`
}

type AppMessageUnreadCountRes struct {
	SystemUnread  uint64 `json:"systemUnread" dc:"系统消息未读数"`
	PrivateUnread uint64 `json:"privateUnread" dc:"私信未读数"`
}
