package gameevent

import "xr-game-server/core/event"

const (
	// RankListRefreshEvent 榜单配置变更(如用户上下榜),触发富豪榜/主播榜刷新
	RankListRefreshEvent event.Type = "RankListRefreshEvent"
)
