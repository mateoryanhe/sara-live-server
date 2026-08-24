package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetMyGuildProfileReq struct {
	g.Meta `path:"/getMyGuildProfile" method:"post" summary:"获取当前CMS用户管理的工会列表" tags:"直播工会"`
}

type MyGuildProfileItem struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	UnsettledTotalIncome float64 `json:"unsettledTotalIncome" dc:"未结算总流水"`
	UpdatedAt            string  `json:"updatedAt"`
}

type GetMyGuildProfileRes struct {
	List []*MyGuildProfileItem `json:"list" dc:"当前CMS用户作为会长管理的工会列表"`
}

type UpdateMyGuildProfileReq struct {
	g.Meta      `path:"/updateMyGuildProfile" method:"post" summary:"更新当前CMS用户管理的指定工会基础信息" tags:"直播工会"`
	ID          uint64 `json:"id" v:"required#工会ID不能为空" dc:"工会ID"`
	Name        string `json:"name" v:"required#工会名称不能为空" dc:"工会名称"`
	Description string `json:"description" dc:"工会简介"`
}

type UpdateMyGuildProfileRes struct {
	Success bool `json:"success"`
}
