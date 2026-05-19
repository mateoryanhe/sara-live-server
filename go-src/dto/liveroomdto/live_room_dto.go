package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// CreateLiveRoomReq 创建直播间(主播ID由鉴权中间件提供)
type CreateLiveRoomReq struct {
	g.Meta `path:"/create" method:"post" summary:"创建直播间" tags:"直播间"`
	Title  string `json:"title"  v:"required|length:1,128#标题不能为空|标题长度需在1到128之间" dc:"直播间标题"`
	Cover  string `json:"cover"  dc:"封面图URL"`
	Notice string `json:"notice" dc:"公告"`
}

type CreateLiveRoomRes struct {
	RoomId  string `json:"roomId"  dc:"直播间ID"`
	GuildId string `json:"guildId" dc:"所属工会ID"`
	Status  uint8  `json:"status"  dc:"状态(0未开播,1直播中)"`
}

// StartLiveReq 开播(主播自身)
type StartLiveReq struct {
	g.Meta `path:"/startLive" method:"post" summary:"开播" tags:"直播间"`
}

type StartLiveRes struct {
	Status uint8 `json:"status" dc:"状态(0未开播,1直播中)"`
}

// StopLiveReq 下播(主播自身)
type StopLiveReq struct {
	g.Meta `path:"/stopLive" method:"post" summary:"下播" tags:"直播间"`
}

type StopLiveRes struct {
	Status uint8 `json:"status" dc:"状态(0未开播,1直播中)"`
}

// UpdateCoverReq 修改封面(主播自身)
type UpdateCoverReq struct {
	g.Meta `path:"/updateCover" method:"post" summary:"修改直播间封面" tags:"直播间"`
	Cover  string `json:"cover" v:"max-length:255#封面URL最长255字符" dc:"封面图URL(允许空字符串清空)"`
}

type UpdateCoverRes struct {
	Success bool `json:"success"`
}

// UpdateNoticeReq 修改公告(主播自身)
type UpdateNoticeReq struct {
	g.Meta `path:"/updateNotice" method:"post" summary:"修改直播间公告" tags:"直播间"`
	Notice string `json:"notice" v:"max-length:512#公告最长512字符" dc:"公告内容(允许空字符串清空)"`
}

type UpdateNoticeRes struct {
	Success bool `json:"success"`
}

// JoinRoomReq 加入直播间(在线玩家登记)
type JoinRoomReq struct {
	g.Meta `path:"/join" method:"post" summary:"加入直播间" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID"`
}

type JoinRoomRes struct {
	OnlineId    string `json:"onlineId"    dc:"在线记录ID(userId_roomId)"`
	OnlineCount int    `json:"onlineCount" dc:"当前在线人数"`
}

// LeaveRoomReq 离开直播间
type LeaveRoomReq struct {
	g.Meta `path:"/leave" method:"post" summary:"离开直播间" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID"`
}

type LeaveRoomRes struct {
	OnlineCount int `json:"onlineCount" dc:"当前在线人数"`
}

// GetOnlineUserListReq 分页查询直播间在线玩家列表
type GetOnlineUserListReq struct {
	g.Meta   `path:"/onlineList" method:"post" summary:"分页查询直播间在线玩家" tags:"直播间"`
	RoomId   uint64 `json:"roomId"   v:"required#直播间ID不能为空" dc:"直播间ID"`
	Page     int    `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int    `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// OnlineUserItem 在线玩家条目
type OnlineUserItem struct {
	UserId   string `json:"userId"   dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Avatar   string `json:"avatar"   dc:"头像URL(已拼资源域名)"`
	JoinedAt int64  `json:"joinedAt" dc:"最近一次加入时间(秒)"`
}

type GetOnlineUserListRes struct {
	Total    int               `json:"total"    dc:"在线总人数"`
	Page     int               `json:"page"     dc:"当前页码"`
	PageSize int               `json:"pageSize" dc:"每页数量"`
	List     []*OnlineUserItem `json:"list"     dc:"玩家列表"`
}

// SendGiftReq App端送礼请求
type SendGiftReq struct {
	g.Meta `path:"/sendGift" method:"post" summary:"直播间送礼" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID"`
	GiftId uint64 `json:"giftId" v:"required#礼物ID不能为空" dc:"礼物ID"`
	Count  int    `json:"count"  v:"required|min:1|max:9999#送礼数量不合法|送礼数量需大于0|单次最多9999" dc:"送礼数量"`
}

// SendGiftRes App端送礼响应
type SendGiftRes struct {
	Cost    uint64  `json:"cost"    dc:"本次实际消耗钻石数"`
	Diamond float64 `json:"diamond" dc:"剩余钻石余额"`
}

// GiftPushItem 推送给房间在线用户的送礼广播载荷
type GiftPushItem struct {
	RoomId       string `json:"roomId"      dc:"直播间ID"`
	SenderId     string `json:"senderId"    dc:"送礼用户ID"`
	SenderName   string `json:"senderName"  dc:"送礼用户昵称"`
	SenderAvatar string `json:"senderAvatar" dc:"送礼用户头像"`
	GiftId       string `json:"giftId"      dc:"礼物ID"`
	GiftName     string `json:"giftName"    dc:"礼物名称"`
	GiftIcon     string `json:"giftIcon"    dc:"礼物图标"`
	GiftAnim     string `json:"giftAnim"    dc:"礼物动画"`
	UnitPrice    uint64 `json:"unitPrice"   dc:"礼物单价"`
	Count        int    `json:"count"       dc:"赠送数量"`
	TotalCost    uint64 `json:"totalCost"   dc:"总消耗钻石数"`
	SentAt       int64  `json:"sentAt"      dc:"发送时间(秒)"`
}

// SendChatReq App端向直播间发送文字消息
type SendChatReq struct {
	g.Meta  `path:"/sendChat" method:"post" summary:"直播间文字消息" tags:"直播间"`
	RoomId  uint64 `json:"roomId"  v:"required#直播间ID不能为空" dc:"直播间ID"`
	Content string `json:"content" v:"required|length:1,256#消息不能为空|消息长度需在1到256之间" dc:"文字内容"`
}

// SendChatRes App端发送文字消息响应
type SendChatRes struct {
	Success bool `json:"success"`
}

// ChatPushItem 推送给房间在线用户的文字消息载荷
type ChatPushItem struct {
	RoomId       string `json:"roomId"       dc:"直播间ID"`
	SenderId     string `json:"senderId"     dc:"发送用户ID"`
	SenderName   string `json:"senderName"   dc:"发送用户昵称"`
	SenderAvatar string `json:"senderAvatar" dc:"发送用户头像"`
	Content      string `json:"content"      dc:"文字内容"`
	SentAt       int64  `json:"sentAt"       dc:"发送时间(秒)"`
}

// GetLiveRoomReq 查询直播间(公开)
type GetLiveRoomReq struct {
	g.Meta `path:"/get" method:"post" summary:"查询直播间" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID"`
}

type GetLiveRoomRes struct {
	RoomId   string `json:"roomId"   dc:"直播间ID(同主播用户ID)"`
	GuildId  string `json:"guildId"  dc:"所属工会ID"`
	Title    string `json:"title"    dc:"直播间标题"`
	Cover    string `json:"cover"    dc:"封面图URL"`
	Notice   string `json:"notice"   dc:"公告"`
	Status   uint8  `json:"status"   dc:"状态(0未开播,1直播中)"`
	CreateAt int64  `json:"createAt" dc:"创建时间(秒)"`
}
