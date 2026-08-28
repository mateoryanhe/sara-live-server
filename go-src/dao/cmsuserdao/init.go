package cmsuserdao

import (
	"strings"

	"xr-game-server/core/cache"
	"xr-game-server/core/httpserver"
	"xr-game-server/entity/cms"
)

func InitCMSUser() {
	cmsLoginUserCacheMgr = cache.NewRowCache[*entity.CMSUser]()
	cmsLoginUserMissCacheMgr = cache.NewRowCache[bool]()
	initRolePermissionCache()
	httpserver.SetCmsApiPermissionChecker(CheckCmsApiPermission)
}

func isDataSyncApiPath(apiPath string) bool {
	return strings.HasPrefix(apiPath, "/dataSync/sync")
}
