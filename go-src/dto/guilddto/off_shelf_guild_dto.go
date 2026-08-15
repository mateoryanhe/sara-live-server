package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// OffShelfGuildListReq CMS工会垃圾库:查询已下架工会(直查DB)
type OffShelfGuildListReq struct {
	g.Meta `path:"/offShelfGuildList" method:"post" summary:"获取下架工会列表(垃圾库)" tags:"直播工会"`
	httpserver.CMSQueryReq
	Name string `json:"name" dc:"工会名称"`
}

// OnShelfGuildReq CMS上架工会
type OnShelfGuildReq struct {
	g.Meta `path:"/onShelfGuild" method:"post" summary:"上架直播工会" tags:"直播工会"`
	ID     uint64 `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
}

type OnShelfGuildRes struct {
	Success bool `json:"success"`
}
