package cfg

import (
	"strings"
	"sync"
)

var (
	runtimeStaticMu    sync.RWMutex
	runtimeStaticPaths []*StaticPathCfg
	runtimeDomainSites []*DomainSiteCfg
	// runtimeStaticSites 按 URL 前缀覆盖 config.yaml 中的一整项静态站点配置。
	// 它与上传资源模块的临时路径/域名分开保存，避免刷新上传配置时误删 H5 等站点映射。
	runtimeStaticSites = make(map[string]*StaticSiteCfg)
)

// RegisterRuntimeStaticPath 注册运行时静态路径(如 CMS 上传资源配置)
func RegisterRuntimeStaticPath(prefix, path string, ttlMinutes int) {
	prefix = normalizeURLPrefix(strings.TrimSpace(prefix))
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	if prefix == "" || prefix == "/" || path == "" {
		return
	}
	runtimeStaticMu.Lock()
	defer runtimeStaticMu.Unlock()
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
	runtimeStaticMu.Lock()
	defer runtimeStaticMu.Unlock()
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

// RegisterRuntimeStaticSite 按 prefix 注册或替换完整静态站点映射。
// 保存后 GetStaticPathRoot、GetStaticSiteDomain 与域名注册表都会读取这一份进程缓存。
func RegisterRuntimeStaticSite(prefix, domain, path string) {
	prefix = normalizeURLPrefix(strings.TrimSpace(prefix))
	domain = strings.TrimSpace(domain)
	path = strings.TrimSpace(path)
	if prefix == "" || prefix == "/" || domain == "" || path == "" {
		return
	}
	runtimeStaticMu.Lock()
	site := cloneStaticSiteCfg(runtimeStaticSites[prefix])
	if site == nil {
		for _, item := range staticSiteCfgs {
			if item != nil && item.Prefix == prefix {
				site = cloneStaticSiteCfg(item)
				break
			}
		}
	}
	if site == nil {
		site = &StaticSiteCfg{}
	}
	site.Domain = domain
	site.Prefix = prefix
	site.Path = path
	site.Root = path
	runtimeStaticSites[prefix] = site
	runtimeStaticMu.Unlock()
}

// UnregisterRuntimeStaticSite 删除指定 prefix 的运行时覆盖，恢复 config.yaml 配置。
func UnregisterRuntimeStaticSite(prefix string) {
	prefix = normalizeURLPrefix(strings.TrimSpace(prefix))
	runtimeStaticMu.Lock()
	delete(runtimeStaticSites, prefix)
	runtimeStaticMu.Unlock()
}

func getRuntimeStaticSite(prefix string) *StaticSiteCfg {
	prefix = normalizeURLPrefix(strings.TrimSpace(prefix))
	runtimeStaticMu.RLock()
	item := cloneStaticSiteCfg(runtimeStaticSites[prefix])
	runtimeStaticMu.RUnlock()
	return item
}

// ClearRuntimeStaticMappings 清空运行时静态映射(CMS 保存前重置)
func ClearRuntimeStaticMappings() {
	runtimeStaticMu.Lock()
	runtimeStaticPaths = nil
	runtimeDomainSites = nil
	runtimeStaticMu.Unlock()
}

func mergedStaticPathCfgs() []*StaticPathCfg {
	sites := mergedStaticSiteCfgs()
	list := buildStaticPathCfgsFromSites(sites)
	runtimeStaticMu.RLock()
	paths := cloneStaticPathCfgs(runtimeStaticPaths)
	runtimeStaticMu.RUnlock()
	for _, override := range paths {
		if override == nil {
			continue
		}
		filtered := list[:0]
		for _, item := range list {
			if item == nil || item.Prefix != override.Prefix {
				filtered = append(filtered, item)
			}
		}
		list = append(filtered, override)
	}
	return normalizeStaticPathCfgs(list)
}

func mergedDomainSiteCfgs() []*DomainSiteCfg {
	list := buildDomainSiteCfgsFromSites(mergedStaticSiteCfgs())
	runtimeStaticMu.RLock()
	domains := cloneDomainSiteCfgs(runtimeDomainSites)
	runtimeStaticMu.RUnlock()
	list = append(list, domains...)
	return normalizeDomainSiteCfgs(list)
}

func mergedStaticSiteCfgs() []*StaticSiteCfg {
	runtimeStaticMu.RLock()
	overrides := make(map[string]*StaticSiteCfg, len(runtimeStaticSites))
	for prefix, item := range runtimeStaticSites {
		overrides[prefix] = cloneStaticSiteCfg(item)
	}
	runtimeStaticMu.RUnlock()

	list := make([]*StaticSiteCfg, 0, len(staticSiteCfgs)+len(overrides))
	for _, item := range staticSiteCfgs {
		if item == nil {
			continue
		}
		if _, replaced := overrides[item.Prefix]; replaced {
			continue
		}
		list = append(list, cloneStaticSiteCfg(item))
	}
	for _, item := range overrides {
		list = append(list, item)
	}
	return normalizeStaticSiteCfgs(list)
}

func cloneStaticSiteCfg(item *StaticSiteCfg) *StaticSiteCfg {
	if item == nil {
		return nil
	}
	cloned := *item
	return &cloned
}

func cloneStaticPathCfgs(list []*StaticPathCfg) []*StaticPathCfg {
	ret := make([]*StaticPathCfg, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		cloned := *item
		ret = append(ret, &cloned)
	}
	return ret
}

func cloneDomainSiteCfgs(list []*DomainSiteCfg) []*DomainSiteCfg {
	ret := make([]*DomainSiteCfg, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		cloned := *item
		ret = append(ret, &cloned)
	}
	return ret
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
