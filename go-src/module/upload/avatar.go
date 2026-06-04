package upload

import "strings"

// GetDefaultAvatarUrl 默认头像完整 URL(未配置时使用内置默认值)
func GetDefaultAvatarUrl() string {
	return getResourceCfgCache().DefaultAvatarUrl
}

// ResolveAvatarUrl 头像文件名转 URL;为空时返回默认头像
func ResolveAvatarUrl(avatarName string) string {
	if strings.TrimSpace(avatarName) == "" {
		return GetDefaultAvatarUrl()
	}
	return GetUrlByName(avatarName)
}
