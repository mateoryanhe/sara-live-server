package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/cms"
)

type CMSSessionModulesReq struct {
	g.Meta `path:"/cmsSessionModules" method:"post" summary:"获取当前CMS用户权限模块" tags:"权限"`
}

type CMSSessionModulesRes struct {
	Admin      bool                 `json:"admin" summary:"是否管理员(含普通/超级)"`
	SuperAdmin bool                 `json:"superAdmin" summary:"是否超级管理员"`
	Modules    []*entity.Permission `json:"modules" summary:"模块列表"`
}
