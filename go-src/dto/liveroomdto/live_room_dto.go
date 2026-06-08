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
}

// StartLiveReq 开播(主播自身)
type StartLiveReq struct {
	g.Meta `path:"/startLive" method:"post" summary:"开播" tags:"直播间"`
}

type StartLiveRes struct {
}

// StopLiveReq 下播(主播自身)
type StopLiveReq struct {
	g.Meta `path:"/stopLive" method:"post" summary:"下播" tags:"直播间"`
}

type StopLiveRes struct {
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
	JoinedAt string `json:"joinedAt" dc:"最近一次加入时间(秒)"`
	Muted    bool   `json:"muted"    dc:"是否被禁言"`
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
	RoomId       uint64 `json:"roomId,string"      dc:"直播间ID"`
	SenderId     uint64 `json:"senderId,string"    dc:"送礼用户ID"`
	SenderName   string `json:"senderName"  dc:"送礼用户昵称"`
	SenderAvatar string `json:"senderAvatar" dc:"送礼用户头像"`
	GiftId       uint64 `json:"giftId,string"      dc:"礼物ID"`
	GiftName     string `json:"giftName"    dc:"礼物名称"`
	GiftIcon     string `json:"giftIcon"    dc:"礼物图标"`
	GiftAnim     string `json:"giftAnim"    dc:"礼物动画"`
	UnitPrice    uint64 `json:"unitPrice"   dc:"礼物单价"`
	Count        int    `json:"count"       dc:"赠送数量"`
	TotalCost    uint64 `json:"totalCost"   dc:"总消耗钻石数"`
	SentAt       int64  `json:"sentAt"      dc:"发送时间(秒)"`
}

// SetAudienceMuteReq 主播对指定观众禁言/解禁
type SetAudienceMuteReq struct {
	g.Meta `path:"/setAudienceMute" method:"post" summary:"主播对观众禁言/解禁" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID(须为本主播房间)"`
	UserId uint64 `json:"userId" v:"required#观众用户ID不能为空" dc:"被操作的观众用户ID"`
	Muted  bool   `json:"muted"  dc:"true禁言,false解禁"`
}

type SetAudienceMuteRes struct {
	Success bool `json:"success"`
}

// AudienceMutePushItem 观众禁言状态推送载荷
type AudienceMutePushItem struct {
	RoomId string `json:"roomId" dc:"直播间ID"`
	UserId string `json:"userId" dc:"观众用户ID"`
	Muted  bool   `json:"muted"  dc:"是否被禁言"`
}

// CancelAudienceMuteReq 主播取消观众禁言
type CancelAudienceMuteReq struct {
	g.Meta `path:"/cancelAudienceMute" method:"post" summary:"主播取消观众禁言" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID(须为本主播房间)"`
	UserId uint64 `json:"userId" v:"required#观众用户ID不能为空" dc:"被取消禁言的观众用户ID"`
}

type CancelAudienceMuteRes struct {
	Success bool `json:"success"`
}

// GetAudienceRestrictStatusReq 查询观众禁言与被踢状态
type GetAudienceRestrictStatusReq struct {
	g.Meta `path:"/getAudienceRestrictStatus" method:"post" summary:"查询观众禁言与被踢状态" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID"`
	UserId uint64 `json:"userId" v:"required#观众用户ID不能为空" dc:"观众用户ID(主播可查任意观众,观众仅可查自己)"`
}

type GetAudienceRestrictStatusRes struct {
	RoomId            string `json:"roomId"            dc:"直播间ID"`
	UserId            string `json:"userId"            dc:"观众用户ID"`
	Muted             bool   `json:"muted"             dc:"是否被禁言"`
	KickBanned        bool   `json:"kickBanned"        dc:"是否处于踢出封禁期"`
	KickTime          int64  `json:"kickTime"          dc:"最近踢出时间(秒,0表示无)"`
	KickBanExpireAt   int64  `json:"kickBanExpireAt"   dc:"踢出封禁截止时间(秒,0表示无)"`
	KickRemainSeconds int64  `json:"kickRemainSeconds" dc:"踢出封禁剩余秒数"`
}

// KickAudienceReq 主播踢出指定观众
type KickAudienceReq struct {
	g.Meta `path:"/kickAudience" method:"post" summary:"主播踢出观众" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID(须为本主播房间)"`
	UserId uint64 `json:"userId" v:"required#观众用户ID不能为空" dc:"被踢出的观众用户ID"`
}

type KickAudienceRes struct {
	Success bool `json:"success"`
}

// AudienceKickPushItem 观众被踢出推送载荷
type AudienceKickPushItem struct {
	RoomId     string `json:"roomId"     dc:"直播间ID"`
	UserId     string `json:"userId"     dc:"观众用户ID"`
	KickTime   int64  `json:"kickTime"   dc:"踢出时间(秒)"`
	BanSeconds int64  `json:"banSeconds" dc:"禁止再次进入时长(秒)"`
}

// CancelKickBanReq 主播取消观众进入限制
type CancelKickBanReq struct {
	g.Meta `path:"/cancelKickBan" method:"post" summary:"主播取消观众进入限制" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required#直播间ID不能为空" dc:"直播间ID(须为本主播房间)"`
	UserId uint64 `json:"userId" v:"required#观众用户ID不能为空" dc:"被取消限制的观众用户ID"`
}

type CancelKickBanRes struct {
	Success bool `json:"success"`
}

// AudienceKickCancelPushItem 取消进入限制推送载荷
type AudienceKickCancelPushItem struct {
	RoomId string `json:"roomId" dc:"直播间ID"`
	UserId string `json:"userId" dc:"观众用户ID"`
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

// AnchorBanPushItem 主播封禁推送载荷(推送给主播及直播间在线观众)
type AnchorBanPushItem struct {
	RoomId       string `json:"roomId"       dc:"直播间ID"`
	AnchorId     string `json:"anchorId"     dc:"主播用户ID"`
	BanApplyTime int64  `json:"banApplyTime" dc:"封禁截止时间(秒)"`
	BannedAt     int64  `json:"bannedAt"     dc:"封禁操作时间(秒)"`
	BanReason    string `json:"banReason"    dc:"封禁原因"`
}
