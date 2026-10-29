package voicechatmicseattierdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type VoiceChatMicSeatTierListReq struct {
	g.Meta `path:"/voiceChatMicSeatTierList" method:"post" summary:"语聊麦位挡位列表" tags:"语聊麦位挡位"`
	httpserver.CMSQueryReq
	SeatCount    int `json:"seatCount" dc:"麦位总数精确匹配,0=不限"`
	StatusFilter int `json:"statusFilter" dc:"0全部,1只看下架,2只看上架"`
}

type VoiceChatMicSeatTierListRes struct {
	ID        string `json:"id"`
	SeatCount int    `json:"seatCount"`
	Sort      int    `json:"sort"`
	Status    uint8  `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateVoiceChatMicSeatTierReq struct {
	g.Meta    `path:"/createVoiceChatMicSeatTier" method:"post" summary:"新建语聊麦位挡位" tags:"语聊麦位挡位"`
	SeatCount int `json:"seatCount" v:"required|min:2|max:25#麦位数量不能为空|麦位数量至少2|麦位数量最多25"`
	Sort      int `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateVoiceChatMicSeatTierRes struct {
	ID string `json:"id"`
}

type UpdateVoiceChatMicSeatTierReq struct {
	g.Meta    `path:"/updateVoiceChatMicSeatTier" method:"post" summary:"修改语聊麦位挡位" tags:"语聊麦位挡位"`
	ID        uint64 `json:"id" v:"required#ID不能为空"`
	SeatCount int    `json:"seatCount" v:"required|min:2|max:25#麦位数量不能为空|麦位数量至少2|麦位数量最多25"`
	Sort      int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateVoiceChatMicSeatTierRes struct {
	Success bool `json:"success"`
}

type DeleteVoiceChatMicSeatTierReq struct {
	g.Meta `path:"/deleteVoiceChatMicSeatTier" method:"post" summary:"删除语聊麦位挡位" tags:"语聊麦位挡位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type DeleteVoiceChatMicSeatTierRes struct {
	Success bool `json:"success"`
}

type OnShelfVoiceChatMicSeatTierReq struct {
	g.Meta `path:"/onShelfVoiceChatMicSeatTier" method:"post" summary:"上架语聊麦位挡位" tags:"语聊麦位挡位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OnShelfVoiceChatMicSeatTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfVoiceChatMicSeatTierReq struct {
	g.Meta `path:"/offShelfVoiceChatMicSeatTier" method:"post" summary:"下架语聊麦位挡位" tags:"语聊麦位挡位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OffShelfVoiceChatMicSeatTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type AppTierListReq struct {
	g.Meta `path:"/appTierList" method:"post" summary:"App查询语聊麦位挡位(已上架)" tags:"语聊麦位挡位"`
}

type AppTierItem struct {
	ID        uint64 `json:"id,string"`
	SeatCount int    `json:"seatCount"`
	Sort      int    `json:"sort"`
}

type AppTierListRes struct {
	List []*AppTierItem `json:"list"`
}
