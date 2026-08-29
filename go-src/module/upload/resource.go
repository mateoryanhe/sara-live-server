package upload

import (
	"net"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// GetResourceDomain 静态资源访问域名
// 本地 debug 未配置 CMS 域名时,自动使用 http://127.0.0.1:端口
func GetResourceDomain() string {
	snap := getResourceCfgCache()
	if snap.ResourceDomainConfigured {
		return snap.ResourceDomain
	}
	if local := resolveDebugLocalResourceDomain(); local != "" {
		return local
	}
	return snap.ResourceDomain
}

// IsResourceDomainConfigured CMS 是否已填写资源域名
func IsResourceDomainConfigured() bool {
	return getResourceCfgCache().ResourceDomainConfigured
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

// buildResourceUrl 拼接 CMS 配置的资源访问根路径 + 相对路径;路径由 CMS 自行填写,服务端不再推导
func buildResourceUrl(name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return name
	}
	base, pathPrefix := parseResourceDomainBase(GetResourceDomain())
	if base == "" {
		return ""
	}
	return joinResourceURL(base, pathPrefix, strings.TrimLeft(name, "/"))
}

func buildImageResourceUrl(fileName string) string {
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return ""
	}
	return buildResourceUrl("/" + fileName)
}

// parseResourceDomainBase 拆分资源域名为 host 根 URL 与可选路径前缀(如 /images)
func parseResourceDomainBase(raw string) (base string, pathPrefix string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return strings.TrimRight(raw, "/"), ""
	}
	base = u.Scheme + "://" + u.Host
	pathPrefix = strings.Trim(strings.ReplaceAll(u.Path, "\\", "/"), "/")
	return base, pathPrefix
}

// joinResourceURL 拼接资源 URL,避免域名路径前缀与文件路径重复
func joinResourceURL(base, pathPrefix, relPath string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	relPath = strings.Trim(strings.ReplaceAll(relPath, "\\", "/"), "/")
	pathPrefix = strings.Trim(strings.ReplaceAll(pathPrefix, "\\", "/"), "/")
	if relPath == "" {
		if pathPrefix == "" {
			return base
		}
		return base + "/" + pathPrefix
	}
	if pathPrefix != "" {
		if relPath == pathPrefix || strings.HasPrefix(relPath, pathPrefix+"/") {
			return base + "/" + relPath
		}
		return base + "/" + pathPrefix + "/" + relPath
	}
	return base + "/" + relPath
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

func isLoopbackHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func resolveDebugLocalResourceDomain() string {
	ctx := gctx.New()
	if !g.Cfg().MustGet(ctx, "server.debugEnv").Bool() {
		return ""
	}
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.address").String())
	if addr == "" {
		return ""
	}
	host := "127.0.0.1"
	port := strings.TrimPrefix(addr, ":")
	if strings.Contains(addr, ":") {
		parts := strings.SplitN(addr, ":", 2)
		if h := strings.TrimSpace(parts[0]); h != "" && h != "0.0.0.0" {
			host = h
		}
		if p := strings.TrimSpace(parts[1]); p != "" {
			port = p
		}
	}
	if port == "" {
		return "http://" + host
	}
	return "http://" + host + ":" + port
}
