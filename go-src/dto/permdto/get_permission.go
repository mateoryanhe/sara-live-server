package permdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PermissionListReq struct {
	g.Meta `path:"/permissionList" method:"post" summary:"获取权限列表" tags:"角色权限"`
	RoleId uint64 `json:"roleId" dc:"角色ID"`
}
