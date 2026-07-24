package upload

import (
	"strings"
)

// GetResourceDomain 静态资源访问域名;未配置时默认 http://127.0.0.1
func GetResourceDomain() string {
	return getResourceCfgCache().ResourceDomain
}

// GetAppImageMaxSize App 端图片上传大小上限(字节);未配置时默认 1MB
func GetAppImageMaxSize() uint64 {
	return getResourceCfgCache().AppImageMaxSize
}

// buildResourceUrl 给资源路径拼接域名;name 为空返回空;已是完整 URL 则原样返回
func buildResourceUrl(name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return name
	}
	domain := GetResourceDomain()
	path := strings.TrimLeft(name, "/")
	return domain + "/" + path
}
