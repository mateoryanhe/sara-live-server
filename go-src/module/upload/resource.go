package upload

import (
	"strings"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/module/globalcfg"
)

// GetResourceDomain 获取静态资源域名(从 GlobalCfg Resource.Domain 读取),未配置返回空串
func GetResourceDomain() string {
	return globalcfgdao.GetStr(globalcfg.Resource, globalcfg.ResourceKeyDomain, "")
}

// BuildResourceUrl 给资源文件名拼接域名,若 name 为空则返回空,若域名未配置则原样返回 name
func buildResourceUrl(name string) string {
	if name == "" {
		return ""
	}
	// 已经是完整 URL,直接返回
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return name
	}
	domain := GetResourceDomain()
	if domain == "" {
		return name
	}
	return strings.TrimRight(domain, "/") + "/" + strings.TrimLeft(name, "/")
}
