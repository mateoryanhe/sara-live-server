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
