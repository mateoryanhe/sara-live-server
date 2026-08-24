package cmsroledto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type RoleListReq struct {
	g.Meta `path:"/roleList" method:"post" summary:"获取角色列表" tags:"角色权限"`
	httpserver.CMSQueryReq
	Name     string `json:"name" dc:"角色名称"`
	RoleType uint8  `json:"roleType" dc:"角色类型 1内部 2外部"`
}

type RoleListRes struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RoleType    uint8  `json:"roleType"`
	Status      uint8  `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
