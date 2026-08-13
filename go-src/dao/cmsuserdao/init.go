package cmsuserdao

import (
	"xr-game-server/core/cache"
	"xr-game-server/core/httpserver"
)

func InitCMSUser() {
	cmsLoginUserCacheMgr = cache.NewCacheMgr()
	cmsLoginUserMissCacheMgr = cache.NewCacheMgr()
	initRolePermissionCache()
	httpserver.SetCmsApiPermissionChecker(CheckCmsApiPermission)
}
