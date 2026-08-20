package game

import (
	"strings"

	"xr-game-server/entity/game"
)

// ResolveAppGameName App 展示名称: 优先直播游戏名称, 为空时回退英文名称.
func ResolveAppGameName(row *entity.GameCfg) string {
	if row == nil {
		return ""
	}
	if name := strings.TrimSpace(row.LiveGameName); name != "" {
		return name
	}
	return strings.TrimSpace(row.NameEn)
}

// ResolveAppGameCover App 展示封面: 优先直播游戏封面, 为空时回退默认封面.
func ResolveAppGameCover(row *entity.GameCfg) string {
	if row == nil {
		return ""
	}
	cover := strings.TrimSpace(row.LiveGameCover)
	if cover == "" {
		cover = strings.TrimSpace(row.Cover)
	}
	return BuildGameCoverUrl(cover)
}
