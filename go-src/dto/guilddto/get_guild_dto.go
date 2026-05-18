package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetGuildReq App端按ID查询工会信息
type GetGuildReq struct {
	g.Meta `path:"/getGuild" method:"post" summary:"获取直播工会信息" tags:"直播工会"`
	ID     uint64 `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
}

type GetGuildRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LeaderId    string `json:"leaderId"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	Status      uint8  `json:"status"`
}
