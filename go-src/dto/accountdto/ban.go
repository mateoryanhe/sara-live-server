package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type BanReq struct {
	g.Meta       `path:"/ban" method:"post" summary:"封号" tags:"账号"`
	AccountId    uint64     `json:"accountId" dc:"账号id"`
	OpenId       string     `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel      uint       `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
	BanApplyTime *time.Time `json:"banApplyTime" dc:"封禁时间"`
}

// BanAnchorReq CMS封禁主播(含App推送,写入 live_rooms)
type BanAnchorReq struct {
	g.Meta       `path:"/banAnchor" method:"post" summary:"封禁主播" tags:"账号"`
	AccountId    uint64     `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
	OpenId       string     `json:"openId" dc:"登陆openId"`
	Channel      uint       `json:"channel" dc:"登陆渠道"`
	BanApplyTime *time.Time `json:"banApplyTime" v:"required#封禁截止时间不能为空" dc:"封禁截止时间"`
	BanReason    string     `json:"banReason" v:"required|length:1,512#封禁原因不能为空|封禁原因长度需在1到512之间" dc:"封禁原因"`
}

// UnBanAnchorReq CMS解封主播直播间
type UnBanAnchorReq struct {
	g.Meta    `path:"/unBanAnchor" method:"post" summary:"解封主播" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
	OpenId    string `json:"openId" dc:"登陆openId"`
	Channel   uint   `json:"channel" dc:"登陆渠道"`
}

type UnBanReq struct {
	g.Meta    `path:"/unBan" method:"post" summary:"解封" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
	OpenId    string `json:"openId" v:"required#openId不能为空" dc:"登陆openId"`
	Channel   uint   `json:"channel" v:"required#channel不能为空" dc:"登陆渠道"`
}

type BanRes struct {
}

// SetLiveRoomStatusReq CMS上架/下架主播直播间
type SetLiveRoomStatusReq struct {
	g.Meta   `path:"/setLiveRoomStatus" method:"post" summary:"上架/下架主播直播间" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播账号ID"`
	Status   uint8  `json:"status" v:"required|in:0,1#状态不能为空|状态只能为0下架或1上架" dc:"状态(0-下架,1-上架)"`
}

type SetLiveRoomStatusRes struct {
	Success bool  `json:"success" dc:"是否成功"`
	Status  uint8 `json:"status" dc:"当前状态(0-下架,1-上架)"`
}

// QueryOffShelfLiveRoomListReq CMS回收站:查询已下架直播间(直查DB)
type QueryOffShelfLiveRoomListReq struct {
	g.Meta `path:"/getOffShelfLiveRoomList" method:"post" summary:"获取下架直播间列表(回收站)" tags:"账号"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"查询关键字(用户ID模糊/昵称/手机号/分享码)"`
}

// OffShelfLiveRoomItem 回收站列表项
type OffShelfLiveRoomItem struct {
	ID           uint64     `json:"id,string"`
	Nickname     string     `json:"nickname"`
	Phone        string     `json:"phone"`
	Avatar       string     `json:"avatar"`
	UserType     uint8      `json:"userType"`
	GuildId      uint64     `json:"guildId,string"`
	RoomTitle    string     `json:"roomTitle"`
	RoomId       uint64     `json:"roomId,string"`
	Category     uint8      `json:"category"`
	Ban          bool       `json:"ban"`
	BanApplyTime *time.Time `json:"banApplyTime"`
	BanReason    string     `json:"banReason"`
	Status       uint8      `json:"status"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	CreatedAt    *time.Time `json:"createdAt"`
}

// ExitGuildReq CMS主播退出工会(将工会ID置为0)
type ExitGuildReq struct {
	g.Meta   `path:"/exitGuild" method:"post" summary:"退出工会" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播账号ID"`
}

type ExitGuildRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
