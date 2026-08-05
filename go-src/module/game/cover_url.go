package game

import (
	"strings"

	"xr-game-server/dao/cfgdao"
)

const vendorGameCoverPathPrefix = "uploads/file/game"

// BuildGameCoverUrl 拼接游戏封面完整 URL(第三方返回的 cover 多为相对路径).
func BuildGameCoverUrl(cover string) string {
	cover = normalizeVendorGameCover(cover)
	if cover == "" {
		return ""
	}
	if strings.HasPrefix(cover, "http://") || strings.HasPrefix(cover, "https://") {
		return cover
	}
	base := cfgdao.GetGameIconBaseUrlFromMemory()
	if base == "" {
		return cover
	}
	if strings.HasPrefix(cover, "/") {
		return base + cover
	}
	return base + "/" + cover
}

// normalizeVendorGameCover 去掉第三方 cover 中的 uploads/file/game 前缀.
func normalizeVendorGameCover(cover string) string {
	cover = strings.TrimSpace(cover)
	if cover == "" {
		return ""
	}
	segment := "/" + vendorGameCoverPathPrefix + "/"
	if idx := strings.Index(strings.ToLower(cover), segment); idx >= 0 {
		cover = cover[idx+len(segment):]
	} else {
		cover = strings.TrimLeft(cover, "/")
		prefix := vendorGameCoverPathPrefix + "/"
		if len(cover) >= len(prefix) && strings.EqualFold(cover[:len(prefix)], prefix) {
			cover = cover[len(prefix):]
		}
	}
	return strings.TrimLeft(cover, "/")
}
