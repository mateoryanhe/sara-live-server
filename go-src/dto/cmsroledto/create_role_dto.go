package cmsroledto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateRoleReq struct {
	g.Meta      `path:"/createRole" method:"post" summary:"创建角色" tags:"角色权限"`
	Name        string   `json:"name" v:"required#角色名称不能为空" dc:"角色名称"`
	Description string   `json:"description" dc:"角色描述"`
	Status      uint8    `json:"status" dc:"状态(0-禁用,1-启用)"`
	Permissions []uint64 `json:"permissions" dc:"权限ID列表"`
}

type CreateRoleRes struct {
	ID uint64 `json:"id"`
}
