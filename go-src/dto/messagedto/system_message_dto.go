package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppSystemMessageUnreadListReq App端查询系统消息未读列表
type AppSystemMessageUnreadListReq struct {
	g.Meta `path:"/systemMessageUnreadList" method:"post" summary:"查询系统消息未读列表" tags:"系统消息"`
}

type AppSystemMessageUnreadListItem struct {
	MsgType     uint8  `json:"msgType" dc:"系统消息类型(1活动消息,2个人系统消息)"`
	UnreadCount uint64 `json:"unreadCount" dc:"未读数"`
	UpdatedAt   string `json:"updatedAt" dc:"更新时间"`
}

type AppSystemMessageUnreadListRes struct {
	List []*AppSystemMessageUnreadListItem `json:"list" dc:"系统消息未读列表"`
}

// AppClearSystemMessageUnreadReq App端上报系统消息已读
type AppClearSystemMessageUnreadReq struct {
	g.Meta    `path:"/clearSystemMessageUnread" method:"post" summary:"上报系统消息已读" tags:"系统消息"`
	MsgType   uint8  `json:"msgType" v:"required|in:1,2#系统消息类型不能为空|系统消息类型无效" dc:"系统消息类型(1活动消息,2个人系统消息)"`
	ReadCount uint64 `json:"readCount" dc:"已读数,为0时表示该类型未读数清零"`
}

type AppClearSystemMessageUnreadRes struct {
	Success      bool   `json:"success" dc:"是否成功"`
	UnreadCount  uint64 `json:"unreadCount" dc:"该类型剩余未读数"`
	SystemUnread uint64 `json:"systemUnread" dc:"系统消息总未读数"`
}

// SystemMessagePushItem 系统消息推送载荷
type SystemMessagePushItem struct {
	Id         uint64 `json:"id,string" dc:"消息ID"`
	ReceiverId uint64 `json:"receiverId,string" dc:"接收者用户ID"`
	Content    string `json:"content" dc:"消息内容"`
	SentAt     string `json:"sentAt" dc:"发送时间"`
}
