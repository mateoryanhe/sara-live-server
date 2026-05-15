package cmsuserdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSUserListReq struct {
	g.Meta `path:"/cmsUserList" method:"post" summary:"获取CMS用户列表" tags:"CMS用户管理"`
	httpserver.CMSQueryReq
	Name   string `json:"name" dc:"CMS用户名称"`
	Status uint8  `json:"status" dc:"状态"`
	Admin  bool   `json:"admin" dc:"是否是管理员"`
}

type CMSUserListRes struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Pwd       string `json:"pwd"`
	Status    uint8  `json:"status"`
	Admin     bool   `json:"admin"`
	RoleId    uint64 `json:"roleId"`
	RoleName  string `json:"roleName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
