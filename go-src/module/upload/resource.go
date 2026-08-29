package upload

import (
	"net"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
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
	// 本地 IP: 存储目录相对 serverRoot(如 /images/a.jpg); 正式服域名: 独立 CDN 根路径(如 /a.jpg)
	if usesLocalServerRootStatic() {
		if rel, ok := imageURLPathUnderServerRoot(fileName); ok {
			return "/" + filepath.ToSlash(rel)
		}
	}
	return "/" + fileName
}

// usesLocalServerRootStatic 本地调试走 config.yaml serverRoot; 正式服域名走 CMS 域名静态站点
func usesLocalServerRootStatic() bool {
	if IsResourceDomainConfigured() {
		return isLoopbackHost(resourceDomainHost(GetResourceDomain()))
	}
	return g.Cfg().MustGet(gctx.New(), "server.debugEnv").Bool()
}

// imageURLPathUnderServerRoot 存储目录在 serverRoot 下时,返回相对 URL 路径(如 images/a.jpg)
func imageURLPathUnderServerRoot(fileName string) (string, bool) {
	storageRoot := filepath.Clean(GetStoragePath())
	serverRoot := filepath.Clean(cfg.GetServerRoot())
	if storageRoot == "" || serverRoot == "" {
		return "", false
	}
	rel, err := filepath.Rel(serverRoot, filepath.Join(storageRoot, fileName))
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return rel, true
}

func buildImageResourceUrl(fileName string) string {
	return buildResourceUrl(buildImageResourcePath(fileName))
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
