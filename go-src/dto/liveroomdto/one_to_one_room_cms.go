package liveroomdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSOneToOneRoomListReq struct {
	g.Meta `path:"/list" method:"post" summary:"1v1房间列表" tags:"1v1房间"`
	httpserver.CMSQueryReq
	Key    string `json:"key" dc:"用户ID/昵称/手机号"`
	Status *uint8 `json:"status" dc:"上下架(0下架,1上架,不传全部)"`
}

type CMSOneToOneRoomItem struct {
	UserId         string     `json:"userId" dc:"主播用户ID"`
	Nickname       string     `json:"nickname"`
	Avatar         string     `json:"avatar"`
	Phone          string     `json:"phone"`
	GuildId        string     `json:"guildId"`
	Title  string `json:"title" dc:"1v1房间标题"`
	Cover  string `json:"cover" dc:"1v1封面URL"`
	TagId  string `json:"tagId" dc:"标签ID"`
	LiveRoomStatus uint8      `json:"liveRoomStatus" dc:"直播间上下架(0下架,1上架)"`
	Status         uint8      `json:"status" dc:"1v1上下架(0下架,1上架)"`
	Billing        float64    `json:"billing" dc:"1v1通话每分钟钻石"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

type CMSCreateOneToOneRoomReq struct {
	g.Meta   `path:"/create" method:"post" summary:"开通1v1房间" tags:"1v1房间"`
	UserId   uint64  `json:"userId,string" v:"required#用户ID不能为空" dc:"主播用户ID"`
	Billing  float64 `json:"billing" v:"min:0#通话单价不能小于0" dc:"1v1通话每分钟钻石"`
	Title string `json:"title" dc:"1v1房间标题"`
	Cover string `json:"cover" dc:"1v1封面对象名"`
	TagId uint64 `json:"tagId,string" dc:"标签ID"`
}

type CMSCreateOneToOneRoomRes struct {
	UserId string `json:"userId"`
}

type CMSUpdateOneToOneRoomReq struct {
	g.Meta   `path:"/update" method:"post" summary:"更新1v1房间" tags:"1v1房间"`
	UserId   uint64  `json:"userId,string" v:"required#用户ID不能为空" dc:"主播用户ID"`
	Billing  float64 `json:"billing" v:"min:0#通话单价不能小于0" dc:"1v1通话每分钟钻石"`
	Title string `json:"title" dc:"1v1房间标题"`
	Cover string `json:"cover" dc:"1v1封面对象名"`
	TagId uint64 `json:"tagId,string" dc:"标签ID"`
}

type CMSUpdateOneToOneRoomRes struct {
}

type CMSSetOneToOneRoomStatusReq struct {
	g.Meta `path:"/setStatus" method:"post" summary:"1v1房间上下架" tags:"1v1房间"`
	UserId uint64 `json:"userId,string" v:"required#用户ID不能为空" dc:"主播用户ID"`
	Status uint8  `json:"status" v:"in:0,1#状态不合法" dc:"0下架,1上架"`
}

type CMSSetOneToOneRoomStatusRes struct {
}
