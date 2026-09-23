package cmsvis

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/liveroomdao"
	cmsentity "xr-game-server/entity/cms"
)

// ResolvePlatformAnchorListVisibility 非管理员只看可见性表中的平台主播，管理员和超管看全部。
func ResolvePlatformAnchorListVisibility(ctx context.Context) (visibleAnchorIds []uint64, filterByVisibility bool) {
	return ResolvePlatformAnchorListVisibilityForUser(httpserver.GetAuthId(ctx))
}

// ResolvePlatformAnchorListVisibilityForUser 按 CMS 用户解析可见平台主播。
func ResolvePlatformAnchorListVisibilityForUser(cmsUserId uint64) (visibleAnchorIds []uint64, filterByVisibility bool) {
	user := cmsuserdao.GetCMSUserById(cmsUserId)
	if cmsentity.CMSUserIsAdmin(user) {
		return nil, false
	}
	return liveroomdao.ListVisiblePlatformAnchorIds(cmsUserId), true
}

// PlatformAnchorVisibilityFilter 返回列表查询需要的可见性过滤条件。
func PlatformAnchorVisibilityFilter(ctx context.Context) (anchorIds []uint64, restrict bool, empty bool) {
	return PlatformAnchorVisibilityFilterForUser(httpserver.GetAuthId(ctx))
}

func PlatformAnchorVisibilityFilterForUser(cmsUserId uint64) (anchorIds []uint64, restrict bool, empty bool) {
	ids, filter := ResolvePlatformAnchorListVisibilityForUser(cmsUserId)
	if !filter {
		return nil, false, false
	}
	if len(ids) == 0 {
		return nil, true, true
	}
	return ids, true, false
}
