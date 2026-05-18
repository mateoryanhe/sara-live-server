package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DeleteGuildReq struct {
	g.Meta `path:"/deleteGuild" method:"post" summary:"删除直播工会" tags:"直播工会"`
	ID     uint64 `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
}

type DeleteGuildRes struct {
	Success bool `json:"success"`
}
