package cmsuser

import (
	"time"

	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/entity"
)

func invalidateCmsToken(userId uint64) {
	if userId == 0 {
		return
	}
	xrtoken.InvalidateCmsTokenCache(userId)
	entity.NewCmsToken(userId, "", time.Now().Add(-time.Hour))
}

func invalidateCmsTokensByRoleId(roleId uint64) {
	if roleId == 0 {
		return
	}
	for _, userId := range cmsuserdao.ListCMSUserIdsByRoleId(roleId) {
		invalidateCmsToken(userId)
	}
}
