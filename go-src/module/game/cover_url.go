package game

import (
	"strings"

	"xr-game-server/module/upload"
)

const vendorGameCoverPathPrefix = "uploads/file/game"

// BuildGameCoverUrl 将自有存储对象路径转换为可访问 URL.
// 完整 URL 仅用于 CMS 人工配置的自定义直播封面.
func BuildGameCoverUrl(cover string) string {
	cover = strings.TrimSpace(cover)
	if cover == "" {
		return ""
	}
	if strings.HasPrefix(cover, "http://") || strings.HasPrefix(cover, "https://") {
		return cover
	}
	return upload.GetUrlByName(cover)
}

// normalizeVendorGameCoverPath 去掉第三方 cover 中的 uploads/file/game 前缀.
// 仅用于拼接同步阶段的第三方下载地址，不能用于数据库展示路径.
func normalizeVendorGameCoverPath(cover string) string {
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
