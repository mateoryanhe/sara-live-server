package upload

import "strings"

// GetDefaultAvatarUrl 默认头像完整 URL(从内置默认头像列表随机取一个)
func GetDefaultAvatarUrl() string {
	return defaultAvatarUrl(PickRandomDefaultAvatarName())
}

// ResolveAvatarUrl 头像文件名转 URL;为空时返回随机默认头像
func ResolveAvatarUrl(avatarName string) string {
	if strings.TrimSpace(avatarName) == "" {
		return GetDefaultAvatarUrl()
	}
	return GetUrlByName(avatarName)
}

// ResolveAvatarUrlForUser 头像文件名转 URL;为空时按 userId 从默认头像列表取一个
func ResolveAvatarUrlForUser(userId uint64, avatarName string) string {
	if strings.TrimSpace(avatarName) == "" {
		return defaultAvatarUrl(PickDefaultAvatarName(userId))
	}
	return GetUrlByName(avatarName)
}
