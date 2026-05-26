package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppSystemMessageListReq App端查询系统消息列表
type AppSystemMessageListReq struct {
	g.Meta    `path:"/systemMessageList" method:"post" summary:"查询系统消息列表" tags:"系统消息"`
	PageIndex int `json:"pageIndex" dc:"页码(从1开始)"`
}

type AppSystemMessageItem struct {
	Id         uint64 `json:"id,string"`
	ReceiverId uint64 `json:"receiverId,string"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
}

type AppSystemMessageListRes struct {
	List []*AppSystemMessageItem `json:"list"`
}

// AppClearSystemMessageUnreadReq App端清除系统消息未读(每次1条)
type AppClearSystemMessageUnreadReq struct {
	g.Meta `path:"/clearSystemMessageUnread" method:"post" summary:"清除系统消息未读" tags:"系统消息"`
}

type AppClearSystemMessageUnreadRes struct {
	Success      bool   `json:"success"`
	SystemUnread uint64 `json:"systemUnread" dc:"剩余系统消息未读数"`
}

// SystemMessagePushItem 系统消息推送载荷
type SystemMessagePushItem struct {
	Id         uint64 `json:"id,string"`
	ReceiverId uint64 `json:"receiverId,string"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	SentAt     string `json:"sentAt"`
}
