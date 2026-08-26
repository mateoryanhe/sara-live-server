package upload

import (
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
)

// GetResourceDomain 静态资源访问域名;未配置时默认 http://127.0.0.1
func GetResourceDomain() string {
	return getResourceCfgCache().ResourceDomain
}

// GetAppImageMaxSize App 端图片上传大小上限(字节);未配置时默认 1MB
func GetAppImageMaxSize() uint64 {
	return getResourceCfgCache().AppImageMaxSize
}

// GetAboutSiteUrl About 页面 URL
func GetAboutSiteUrl() string {
	return buildResourceUrl("/about.html")
}

// GetSafetyCenterUrl 安全中心页面 URL
func GetSafetyCenterUrl() string {
	return buildResourceUrl("/safety-center.html")
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

func buildImageResourcePath(fileName string) string {
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return ""
	}
	if resourceDomainUsesDedicatedImageHost(GetResourceDomain()) {
		return "/" + fileName
	}
	segment := cfg.GetImageStaticPathSegment()
	if segment == "" {
		return "/" + fileName
	}
	return "/" + segment + "/" + fileName
}

func buildImageResourceUrl(fileName string) string {
	return buildResourceUrl(buildImageResourcePath(fileName))
}

func resourceDomainUsesDedicatedImageHost(domain string) bool {
	segment := strings.Trim(cfg.GetImageStaticPathSegment(), "/")
	if segment == "" {
		return false
	}
	host := resourceDomainHost(domain)
	if host == "" {
		return false
	}
	host = strings.ToLower(host)
	return host == segment || strings.HasPrefix(host, segment+".")
}

func resourceDomainHost(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return ""
	}
	if u, err := url.Parse(domain); err == nil && u.Host != "" {
		return u.Hostname()
	}
	if u, err := url.Parse("https://" + strings.TrimPrefix(domain, "//")); err == nil && u.Host != "" {
		return u.Hostname()
	}
	return ""
}
