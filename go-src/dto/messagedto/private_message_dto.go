package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppSendPrivateMessageReq App端发送私信
type AppSendPrivateMessageReq struct {
	g.Meta     `path:"/sendPrivateMessage" method:"post" summary:"发送私信" tags:"私信"`
	ReceiverId uint64 `json:"receiverId" v:"required|min:1#接收者ID不能为空|接收者ID非法" dc:"接收者用户ID"`
	Content    string `json:"content" v:"required|length:1,1024#消息内容不能为空|消息内容长度需在1到1024之间" dc:"消息内容"`
}

type AppSendPrivateMessageRes struct {
	MessageId uint64 `json:"messageId,string" dc:"消息ID"`
	Success   bool   `json:"success"`
}

// AppPrivateMessageListReq App端查询收到的私信
type AppPrivateMessageListReq struct {
	g.Meta    `path:"/privateMessageList" method:"post" summary:"查询私信列表" tags:"私信"`
	PageIndex int `json:"pageIndex" dc:"页码(从1开始)"`
}

type AppPrivateMessageItem struct {
	Id           uint64 `json:"id,string"`
	SenderId     uint64 `json:"senderId,string"`
	ReceiverId   uint64 `json:"receiverId,string"`
	Content      string `json:"content"`
	SenderName   string `json:"senderName"`
	SenderAvatar string `json:"senderAvatar"`
	CreatedAt    string `json:"createdAt"`
}

type AppPrivateMessageListRes struct {
	List []*AppPrivateMessageItem `json:"list"`
}

// PrivateMessagePushItem 私信推送载荷
type PrivateMessagePushItem struct {
	Id           uint64 `json:"id,string"`
	SenderId     uint64 `json:"senderId,string"`
	ReceiverId   uint64 `json:"receiverId,string"`
	Content      string `json:"content,string"`
	SenderName   string `json:"senderName"`
	SenderAvatar string `json:"senderAvatar"`
	SentAt       string `json:"sentAt"`
}
