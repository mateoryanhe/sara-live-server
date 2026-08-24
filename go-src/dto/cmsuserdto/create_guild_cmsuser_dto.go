package cmsuserdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateGuildCMSUserReq struct {
	g.Meta `path:"/createGuildCMSUser" method:"post" summary:"工会模块创建外部CMS用户" tags:"CMS用户管理"`
	Name   string `json:"name" v:"required#CMS用户名称不能为空" dc:"CMS用户名称"`
	Pwd    string `json:"pwd" v:"required#密码不能为空" dc:"密码"`
	Status uint8  `json:"status" dc:"状态(0-禁用,1-启用)"`
	RoleId uint64 `json:"roleId" v:"required#角色不能为空" dc:"角色ID"`
	Remark string `json:"remark" dc:"备注"`
}
