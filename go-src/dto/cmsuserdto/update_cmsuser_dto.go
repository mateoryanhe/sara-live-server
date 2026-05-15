package cmsuserdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UpdateCMSUserReq struct {
	g.Meta `path:"/updateCMSUser" method:"post" summary:"更新CMS用户" tags:"CMS用户管理"`
	ID     uint64 `json:"id" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
	Name   string `json:"name" v:"required#CMS用户名称不能为空" dc:"CMS用户名称"`
	Pwd    string `json:"pwd" dc:"密码"`
	Status uint8  `json:"status" dc:"状态(0-禁用,1-启用)"`
	Admin  bool   `json:"admin" dc:"是否是管理员"`
	RoleId uint64 `json:"roleId" dc:"角色ID"`
}

type UpdateCMSUserRes struct {
	Success bool `json:"success"`
}
