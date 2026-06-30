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

// AppPrivateMessageUnreadListReq App端分页查询私信未读明细
type AppPrivateMessageUnreadListReq struct {
	g.Meta    `path:"/privateMessageUnreadList" method:"post" summary:"查询私信未读明细列表" tags:"私信"`
	PageIndex int `json:"pageIndex" dc:"页码(从1开始)"`
}

type AppPrivateMessageUnreadDetailItem struct {
	SenderId     uint64 `json:"sender_id,string"`
	UnreadCount  uint64 `json:"unreadCount"`
	SenderName   string `json:"senderName"`
	SenderAvatar string `json:"senderAvatar"`
	UpdatedAt    string `json:"updatedAt"`
}

type AppPrivateMessageUnreadListRes struct {
	List []*AppPrivateMessageUnreadDetailItem `json:"list"`
}

// AppPrivateMessageBySenderReq App端按发送者查询私信内容
type AppPrivateMessageBySenderReq struct {
	g.Meta        `path:"/privateMessageBySender" method:"post" summary:"按目标用户id查询私信" tags:"私信"`
	TargetId      uint64 `json:"targetId,string" v:"required|min:1#发送者ID不能为空|发送者ID非法" dc:"目标用户ID"`
	LastCreatedAt int64  `json:"lastCreatedAt" dc:"游标(毫秒时间戳);传0从最新消息开始,传上一页最后一条消息的创建时间戳拉取更早消息"`
	PageSize      int    `json:"pageSize" dc:"每页条数(默认40)"`
}

type AppPrivateMessageItem struct {
	Id           uint64 `json:"id,string"`
	SenderId     uint64 `json:"senderId,string"`
	ReceiverId   uint64 `json:"receiverId,string"`
	Content      string `json:"content"`
	SenderName   string `json:"senderName"`
	SenderAvatar string `json:"senderAvatar"`
	CreatedAt    string `json:"createdAt"`
	CreatedAtMs  int64  `json:"createdAtMs" dc:"创建时间(毫秒),用作下一页 lastCreatedAt 游标"`
}

type AppPrivateMessageBySenderRes struct {
	List    []*AppPrivateMessageItem `json:"list"`
	HasMore bool                     `json:"hasMore" dc:"是否还有更早消息"`
}

// AppClearPrivateMessageUnreadReq App端清除指定玩家私信未读
type AppClearPrivateMessageUnreadReq struct {
	g.Meta       `path:"/clearPrivateMessageUnread" method:"post" summary:"清除指定玩家私信未读" tags:"私信"`
	TargetId     uint64 `json:"targetId,string" v:"required|min:1#发送者ID不能为空|发送者ID非法" dc:"目标用户ID"`
	ClearedCount uint64 `json:"clearedCount" dc:"本次清除的未读数"`
}

type AppClearPrivateMessageUnreadRes struct {
	Success       bool   `json:"success"`
	ClearedCount  uint64 `json:"clearedCount" dc:"本次清除的未读数"`
	PrivateUnread uint64 `json:"privateUnread" dc:"剩余私信未读总数"`
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
