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

// SystemMessagePushItem 系统消息推送载荷
type SystemMessagePushItem struct {
	Id         uint64 `json:"id,string"`
	ReceiverId uint64 `json:"receiverId,string"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	SentAt     string `json:"sentAt"`
}
