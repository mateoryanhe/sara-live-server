package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UpdateGuildReq struct {
	g.Meta       `path:"/updateGuild" method:"post" summary:"更新直播工会" tags:"直播工会"`
	ID           uint64   `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
	Name         string   `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	LeaderId     uint64   `json:"leaderId" dc:"会长/负责人ID"`
	Description  string   `json:"description" dc:"工会简介"`
	GuildType    uint8    `json:"guildType" v:"in:0,1#工会类型仅支持0普通或1币商" dc:"工会类型(0普通,1币商)"`
	SharePercent *float64 `json:"sharePercent" dc:"分佣比例(%);仅币商工会可自定义,空则保持原值"`
}

type UpdateGuildRes struct {
	Success bool `json:"success"`
}
