package cmsvis

import (
	"context"

	"xr-game-server/core/httpserver"
)

// AnchorRoomVisibility CMS 主播房间可见性(工会 OR 平台主播,与主播列表 visibilityScoped 一致)
type AnchorRoomVisibility struct {
	Restrict          bool
	Empty             bool
	GuildIds          []uint64
	PlatformAnchorIds []uint64
}

// ResolveAnchorRoomVisibility 解析当前请求 CMS 用户的主播房间可见性
func ResolveAnchorRoomVisibility(ctx context.Context) AnchorRoomVisibility {
	return ResolveAnchorRoomVisibilityForUser(httpserver.GetAuthId(ctx))
}

// ResolveAnchorRoomVisibilityForUser 按 CMS 用户解析主播房间可见性(异步导出等无 request ctx 时使用)
func ResolveAnchorRoomVisibilityForUser(cmsUserId uint64) AnchorRoomVisibility {
	guildIds, guildRestrict, guildEmpty := VisibilityGuildFilterForUser(cmsUserId)
	platformIds, platformRestrict, platformEmpty := PlatformAnchorVisibilityFilterForUser(cmsUserId)
	if !guildRestrict && !platformRestrict {
		return AnchorRoomVisibility{Restrict: false, Empty: false}
	}
	hasGuild := guildRestrict && !guildEmpty && len(guildIds) > 0
	hasPlatform := platformRestrict && !platformEmpty && len(platformIds) > 0
	if !hasGuild && !hasPlatform {
		return AnchorRoomVisibility{Restrict: true, Empty: true}
	}
	vis := AnchorRoomVisibility{Restrict: true, Empty: false}
	if hasGuild {
		vis.GuildIds = guildIds
	}
	if hasPlatform {
		vis.PlatformAnchorIds = platformIds
	}
	return vis
}
