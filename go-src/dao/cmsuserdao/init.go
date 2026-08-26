package cmsuserdao

import (
	"xr-game-server/entity/cms"
	"xr-game-server/core/cache"
	"xr-game-server/core/httpserver"
)

func InitCMSUser() {
	cmsLoginUserCacheMgr = cache.NewRowCache[*entity.CMSUser]()
	cmsLoginUserMissCacheMgr = cache.NewRowCache[bool]()
	initRolePermissionCache()
	httpserver.SetCmsApiPermissionChecker(CheckCmsApiPermission)
}
