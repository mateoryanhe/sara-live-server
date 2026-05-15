package cmsroledto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DeleteRoleReq struct {
	g.Meta `path:"/deleteRole" method:"post" summary:"删除角色" tags:"角色权限"`
	ID     uint64 `json:"id" v:"required#角色ID不能为空" dc:"角色ID"`
}

type DeleteRoleRes struct {
	Success bool `json:"success"`
}
