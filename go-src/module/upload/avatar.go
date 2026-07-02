package upload

import "strings"

// ResolveAvatarUrlForUser 头像文件名转 URL;为空时按 userId 从默认头像列表取一个
func ResolveAvatarUrlForUser(userId uint64, avatarName string) string {
	if strings.TrimSpace(avatarName) == "" {
		return defaultAvatarUrl(PickDefaultAvatarName(userId))
	}
	return GetUrlByName(avatarName)
}
