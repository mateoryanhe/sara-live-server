package cfg

import (
	"strings"
)

var runtimeStaticPaths []*StaticPathCfg
var runtimeDomainSites []*DomainSiteCfg

// RegisterRuntimeStaticPath 注册运行时静态路径(如 CMS 上传资源配置)
func RegisterRuntimeStaticPath(prefix, path string, ttlMinutes int) {
	prefix = normalizeURLPrefix(strings.TrimSpace(prefix))
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	if prefix == "" || prefix == "/" || path == "" {
		return
	}
	for _, item := range runtimeStaticPaths {
		if item != nil && item.Prefix == prefix {
			item.Path = path
			item.Root = path
			item.TTLMinutes = ttlMinutes
			return
		}
	}
	runtimeStaticPaths = append(runtimeStaticPaths, &StaticPathCfg{
		Prefix:     prefix,
		Path:       path,
		Root:       path,
		TTLMinutes: ttlMinutes,
	})
}

// RegisterRuntimeDomainSite 注册运行时域名静态站点
func RegisterRuntimeDomainSite(domain, root string) {
	domain = strings.TrimSpace(domain)
	root = strings.TrimRight(strings.TrimSpace(root), "/")
	if domain == "" || root == "" {
		return
	}
	for _, item := range runtimeDomainSites {
		if item != nil && strings.EqualFold(strings.TrimSpace(item.Domain), domain) {
			item.Root = root
			return
		}
	}
	runtimeDomainSites = append(runtimeDomainSites, &DomainSiteCfg{
		Domain: domain,
		Root:   root,
	})
}

// ClearRuntimeStaticMappings 清空运行时静态映射(CMS 保存前重置)
func ClearRuntimeStaticMappings() {
	runtimeStaticPaths = nil
	runtimeDomainSites = nil
}

func mergedStaticPathCfgs() []*StaticPathCfg {
	if len(runtimeStaticPaths) == 0 {
		return staticPathCfgs
	}
	return normalizeStaticPathCfgs(append(append([]*StaticPathCfg{}, staticPathCfgs...), runtimeStaticPaths...))
}

func mergedDomainSiteCfgs() []*DomainSiteCfg {
	if len(runtimeDomainSites) == 0 {
		return domainSiteCfgs
	}
	return normalizeDomainSiteCfgs(append(append([]*DomainSiteCfg{}, domainSiteCfgs...), runtimeDomainSites...))
}

// RefreshImageStaticCfg CMS 上传资源路径注册后刷新图片静态前缀
func RefreshImageStaticCfg() {
	prefix := normalizeURLPrefix("/images")
	root := strings.TrimSpace(GetStaticPathRoot(prefix))
	if root == "" {
		return
	}
	imageStaticPrefix = prefix
}
