package cmsvis

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	cmsentity "xr-game-server/entity/cms"
)

// ResolveGuildListVisibility 非管理员仅看可见性表中的工会;管理员/超管看全部
func ResolveGuildListVisibility(ctx context.Context) (visibleGuildIds []uint64, filterByVisibility bool) {
	return ResolveGuildListVisibilityForUser(httpserver.GetAuthId(ctx))
}

// ResolveGuildListVisibilityForUser 按 CMS 用户解析可见工会(异步导出等无 request ctx 时使用)
func ResolveGuildListVisibilityForUser(cmsUserId uint64) (visibleGuildIds []uint64, filterByVisibility bool) {
	user := cmsuserdao.GetCMSUserById(cmsUserId)
	if cmsentity.CMSUserIsAdmin(user) {
		return nil, false
	}
	return guilddao.ListVisibleGuildIds(cmsUserId), true
}

// VisibilityGuildFilter 列表/导出共用:restrict=true 时须按 guildIds 过滤;empty=true 表示应直接返回空结果
func VisibilityGuildFilter(ctx context.Context) (guildIds []uint64, restrict bool, empty bool) {
	return VisibilityGuildFilterForUser(httpserver.GetAuthId(ctx))
}

// VisibilityGuildFilterForUser 按 CMS 用户解析可见性过滤条件
func VisibilityGuildFilterForUser(cmsUserId uint64) (guildIds []uint64, restrict bool, empty bool) {
	ids, filter := ResolveGuildListVisibilityForUser(cmsUserId)
	if !filter {
		return nil, false, false
	}
	if len(ids) == 0 {
		return nil, true, true
	}
	return ids, true, false
}
