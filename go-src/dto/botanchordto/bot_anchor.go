package botanchordto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// QueryBotAnchorListReq CMS分页查询机器人主播列表
type QueryBotAnchorListReq struct {
	g.Meta `path:"/getBotAnchorList" method:"post" summary:"获取机器人主播列表" tags:"机器人主播"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"查询关键字(用户ID/昵称)"`
}

// BotAnchorListItem 机器人主播列表项
type BotAnchorListItem struct {
	ID                   uint64     `json:"id,string"`
	Nickname             string     `json:"nickname"`
	Avatar               string     `json:"avatar"`
	GuildId              uint64     `json:"guildId,string"`
	RoomId               uint64     `json:"roomId,string"`
	RoomTitle            string     `json:"roomTitle" dc:"直播间标题"`
	Category             uint8      `json:"category" dc:"直播间类型(1=hot,2=game,3=私密)"`
	TagId                uint64     `json:"tagId,string" dc:"直播间标签ID"`
	TagName              string     `json:"tagName" dc:"直播间标签名称"`
	CloudPlayerVideo     string     `json:"cloudPlayerVideo" dc:"云播放器MP4视频访问URL"`
	CloudPlayerVideoFile string     `json:"cloudPlayerVideoFile" dc:"云播放器MP4视频存储文件名/路径"`
	PushStream           bool       `json:"pushStream" dc:"是否推流"`
	IsTest               bool       `json:"isTest" dc:"是否测试机器人主播"`
	BotAnchorStatus      uint8      `json:"botAnchorStatus" dc:"状态(0停用,1启用)"`
	LiveStatus           uint8      `json:"liveStatus" dc:"直播状态(0未开播,1直播中)"`
	CreatedAt            *time.Time `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
}

// CreateBotAnchorReq CMS创建机器人主播
type CreateBotAnchorReq struct {
	g.Meta           `path:"/createBotAnchor" method:"post" summary:"创建机器人主播" tags:"机器人主播"`
	Nickname         string `json:"nickname" v:"required|length:1,32" dc:"昵称"`
	Avatar           string `json:"avatar" dc:"头像文件名"`
	GuildId          uint64 `json:"guildId,string" dc:"工会ID(可选)"`
	RoomTitle        string `json:"roomTitle" dc:"直播间标题"`
	Category         uint8  `json:"category" v:"in:1,2,3" dc:"直播间类型(1=hot,2=game,3=私密)"`
	TagId            uint64 `json:"tagId,string" dc:"直播间标签ID(0表示无)"`
	CloudPlayerVideo string `json:"cloudPlayerVideo" dc:"云播放器MP4视频文件名或URL"`
	PushStream       bool   `json:"pushStream" dc:"是否推流"`
	IsTest           bool   `json:"isTest" dc:"是否测试机器人主播"`
}

// CreateBotAnchorRes CMS创建机器人主播响应
type CreateBotAnchorRes struct {
	ID uint64 `json:"id,string"`
}

// UpdateBotAnchorReq CMS更新机器人主播资料
type UpdateBotAnchorReq struct {
	g.Meta           `path:"/updateBotAnchor" method:"post" summary:"更新机器人主播" tags:"机器人主播"`
	ID               uint64  `json:"id,string" v:"required" dc:"用户ID"`
	Nickname         string  `json:"nickname" v:"required|length:1,32" dc:"昵称"`
	Avatar           *string `json:"avatar" dc:"头像文件名,不传表示不修改"`
	RoomTitle        string  `json:"roomTitle" dc:"直播间标题"`
	Category         uint8   `json:"category" v:"in:1,2,3" dc:"直播间类型(1=hot,2=game,3=私密)"`
	TagId            uint64  `json:"tagId,string" dc:"直播间标签ID(0表示无)"`
	CloudPlayerVideo *string `json:"cloudPlayerVideo" dc:"云播放器MP4视频文件名或URL,不传表示不修改"`
	PushStream       bool    `json:"pushStream" dc:"是否推流"`
	IsTest           bool    `json:"isTest" dc:"是否测试机器人主播"`
}

// UpdateBotAnchorRes CMS更新机器人主播响应
type UpdateBotAnchorRes struct {
	Success bool `json:"success"`
}

// SetBotAnchorStatusReq CMS启用/停用机器人主播
type SetBotAnchorStatusReq struct {
	g.Meta `path:"/setBotAnchorStatus" method:"post" summary:"设置机器人主播状态" tags:"机器人主播"`
	ID     uint64 `json:"id,string" v:"required" dc:"用户ID"`
	Status uint8  `json:"status" v:"in:0,1" dc:"状态(0停用,1启用)"`
}

// SetBotAnchorStatusRes CMS设置机器人主播状态响应
type SetBotAnchorStatusRes struct {
	Success bool `json:"success"`
}

// StartBotAnchorLiveReq CMS机器人主播开播
type StartBotAnchorLiveReq struct {
	g.Meta `path:"/startBotAnchorLive" method:"post" summary:"机器人主播开播" tags:"机器人主播"`
	ID     uint64 `json:"id,string" v:"required" dc:"用户ID"`
}

// StartBotAnchorLiveRes CMS机器人主播开播响应
type StartBotAnchorLiveRes struct {
	Success bool `json:"success"`
}

// StopBotAnchorLiveReq CMS机器人主播下播
type StopBotAnchorLiveReq struct {
	g.Meta `path:"/stopBotAnchorLive" method:"post" summary:"机器人主播下播" tags:"机器人主播"`
	ID     uint64 `json:"id,string" v:"required" dc:"用户ID"`
}

// StopBotAnchorLiveRes CMS机器人主播下播响应
type StopBotAnchorLiveRes struct {
	Success bool `json:"success"`
}

// BatchStartBotAnchorLiveReq CMS批量机器人主播开播
type BatchStartBotAnchorLiveReq struct {
	g.Meta `path:"/batchStartBotAnchorLive" method:"post" summary:"批量机器人主播开播" tags:"机器人主播"`
	IDs    []uint64 `json:"ids,string" v:"required|min-length:1#请至少选择一个机器人主播"`
}

// BatchStopBotAnchorLiveReq CMS批量机器人主播下播
type BatchStopBotAnchorLiveReq struct {
	g.Meta `path:"/batchStopBotAnchorLive" method:"post" summary:"批量机器人主播下播" tags:"机器人主播"`
	IDs    []uint64 `json:"ids,string" v:"required|min-length:1#请至少选择一个机器人主播"`
}

// BatchBotAnchorLiveRes CMS批量开播/下播响应
type BatchBotAnchorLiveRes struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	FailIds      []uint64 `json:"failIds,string"`
}
