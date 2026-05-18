package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateGuildReq struct {
	g.Meta      `path:"/createGuild" method:"post" summary:"创建直播工会" tags:"直播工会"`
	Name        string `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	LeaderId    uint64 `json:"leaderId" dc:"会长/负责人ID"`
	Contact     string `json:"contact" dc:"联系方式"`
	Description string `json:"description" dc:"工会简介"`
	Status      uint8  `json:"status" dc:"状态(0-禁用,1-启用)"`
}

type CreateGuildRes struct {
	ID string `json:"id"`
}
