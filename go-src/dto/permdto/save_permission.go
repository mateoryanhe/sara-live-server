package permdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

type CreatePermissionReq struct {
	g.Meta `path:"/createPermission" method:"post" summary:"创建权限" tags:"角色权限"`
	Data   []*entity.Permission `json:"data" v:"required#权限数据不能为空" dc:"权限数据"`
}
