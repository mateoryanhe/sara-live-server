package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GuildListReq struct {
	g.Meta `path:"/guildList" method:"post" summary:"获取直播工会列表" tags:"直播工会"`
	httpserver.CMSQueryReq
	Name string `json:"name" dc:"工会名称"`
}

type GuildListRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LeaderId    string `json:"leaderId"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	Status      uint8  `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
