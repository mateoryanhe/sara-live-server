package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UpdateGuildReq struct {
	g.Meta      `path:"/updateGuild" method:"post" summary:"更新直播工会" tags:"直播工会"`
	ID          uint64 `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
	Name        string `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	LeaderId    uint64 `json:"leaderId" dc:"会长/负责人ID"`
	Description string `json:"description" dc:"工会简介"`
}

type UpdateGuildRes struct {
	Success bool `json:"success"`
}
