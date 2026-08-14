package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateGuildReq struct {
	g.Meta      `path:"/createGuild" method:"post" summary:"创建直播工会" tags:"直播工会"`
	Name        string `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	LeaderId    uint64 `json:"leaderId" dc:"会长/负责人ID"`
	Description string `json:"description" dc:"工会简介"`
}

type CreateGuildRes struct {
	ID string `json:"id"`
}
