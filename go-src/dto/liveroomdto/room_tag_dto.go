package liveroomdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type LiveRoomTagListReq struct {
	g.Meta `path:"/liveRoomTagList" method:"post" summary:"获取直播间标签列表" tags:"直播间标签"`
	httpserver.CMSQueryReq
}

type LiveRoomTagListRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Sort      int    `json:"sort"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateLiveRoomTagReq struct {
	g.Meta `path:"/createLiveRoomTag" method:"post" summary:"创建直播间标签" tags:"直播间标签"`
	Name   string `json:"name" v:"required|length:1,64#标签名称不能为空|标签名称长度需在1-64字符" dc:"标签名称"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateLiveRoomTagRes struct {
	ID string `json:"id" dc:"标签ID"`
}

type UpdateLiveRoomTagReq struct {
	g.Meta `path:"/updateLiveRoomTag" method:"post" summary:"修改直播间标签" tags:"直播间标签"`
	ID     uint64 `json:"id" v:"required#标签ID不能为空" dc:"标签ID"`
	Name   string `json:"name" v:"required|length:1,64#标签名称不能为空|标签名称长度需在1-64字符" dc:"标签名称"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateLiveRoomTagRes struct {
	Success bool `json:"success"`
}

type DeleteLiveRoomTagReq struct {
	g.Meta `path:"/deleteLiveRoomTag" method:"post" summary:"删除直播间标签" tags:"直播间标签"`
	ID     uint64 `json:"id" v:"required#标签ID不能为空" dc:"标签ID"`
}

type DeleteLiveRoomTagRes struct {
	Success bool `json:"success"`
}

type AppLiveRoomTagListReq struct {
	g.Meta `path:"/appLiveRoomTagList" method:"post" summary:"App查询直播间标签列表" tags:"直播间标签"`
}

type AppLiveRoomTagItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

type AppLiveRoomTagListRes struct {
	List []*AppLiveRoomTagItem `json:"list"`
}
