package cmsuserdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSUserListReq struct {
	g.Meta `path:"/cmsUserList" method:"post" summary:"获取CMS用户列表" tags:"CMS用户管理"`
	httpserver.CMSQueryReq
	Name     string `json:"name" dc:"CMS用户名称"`
	Key      string `json:"key" dc:"关键字(用户名/ID模糊)"`
	RoleId   string `json:"roleId" dc:"角色ID(可选)"`
	Status   uint8  `json:"status" dc:"状态"`
	Admin     bool   `json:"admin" dc:"是否是管理员(兼容)"`
	AdminType uint8  `json:"adminType" dc:"管理员类型 0全部 1普通管理员 2超级管理员"`
	RoleType  uint8  `json:"roleType" dc:"按关联角色类型筛选 1内部 2外部"`
	NonAdmin bool   `json:"nonAdmin" dc:"仅非管理员"`
}

type CMSUserListRes struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Pwd       string `json:"pwd"`
	Status    uint8  `json:"status"`
	Admin     bool   `json:"admin"`
	AdminType uint8  `json:"adminType"`
	RoleId    uint64 `json:"roleId"`
	RoleName  string `json:"roleName"`
	RoleType  uint8  `json:"roleType"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
