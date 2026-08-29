package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/cfg"
)

type domainSiteEntry struct {
	domains []string
	root    string
}

type domainCertEntry struct {
	domains []string
	cert    tls.Certificate
}

var domainSiteEntries []domainSiteEntry
var domainCertEntries []domainCertEntry

func setupDomainSites() {
	sites := cfg.GetDomainSiteCfgs()
	if len(sites) == 0 {
		return
	}
	ctx := gctx.New()
	domainSiteEntries = make([]domainSiteEntry, 0, len(sites))
	for _, site := range sites {
		root := resolveDomainSiteRoot(ctx, site.Root)
		if root == "" {
			g.Log().Warningf(ctx, "域名站点根目录无效,已跳过 domain=%s root=%s", site.Domain, site.Root)
			continue
		}
		domainSiteEntries = append(domainSiteEntries, domainSiteEntry{
			domains: cfg.SplitDomains(site.Domain),
			root:    root,
		})
		if site.CertFile != "" && site.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(site.CertFile, site.KeyFile)
			if err != nil {
				g.Log().Errorf(ctx, "加载域名证书失败 domain=%s cert=%s key=%s err=%v", site.Domain, site.CertFile, site.KeyFile, err)
				continue
			}
			domainCertEntries = append(domainCertEntries, domainCertEntry{
				domains: cfg.SplitDomains(site.Domain),
				cert:    cert,
			})
		}
	}
	if tlsConfig := buildDomainTLSConfig(ctx); tlsConfig != nil {
		httpServer.SetTLSConfig(tlsConfig)
		g.Log().Warning(ctx, "已启用多域名HTTPS证书(SNI)")
	}
	if len(domainSiteEntries) > 0 {
		g.Log().Warningf(ctx, "已启用域名静态目录映射,共 %d 项", len(domainSiteEntries))
		bindDomainStaticHooks()
	}
}

func resolveDomainSiteRoot(ctx context.Context, root string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		return ""
	}
	if real := gfile.RealPath(root); real != "" {
		return real
	}
	if err := gfile.Mkdir(root); err != nil {
		g.Log().Warningf(ctx, "创建域名站点根目录失败 root=%s err=%v", root, err)
	}
	if real := gfile.RealPath(root); real != "" {
		return real
	}
	return filepath.Clean(root)
}

func bindDomainStaticHooks() {
	for _, entry := range domainSiteEntries {
		entry := entry
		httpServer.Domain(strings.Join(entry.domains, ",")).BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
			serveDomainStatic(r, entry.root)
		})
	}
}

func serveDomainStatic(r *ghttp.Request, root string) {
	if r == nil || root == "" {
		return
	}
	if handleStaticCORS(r) {
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}
	reqPath := mapStaticSiteRequestPath(root, r.URL.Path)
	cmsSite := isCMSDomainRoot(root)
	if cmsSite && reqPath == "" {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	if !cmsSite && (reqPath == "/cms" || strings.HasPrefix(reqPath, "/cms/")) {
		return
	}
	filePath, ok := buildDomainStaticFilePath(root, reqPath)
	if ok && gfile.Exists(filePath) && !gfile.IsDir(filePath) {
		writeStaticFile(r, filePath)
		r.ExitAll()
		return
	}
	if cmsSite {
		serveCMSDomainSPA(r, root, reqPath)
		return
	}
	// 非 CMS 站点且文件不存在时不拦截,交给 serverRoot / API 路由继续处理
}

func mapStaticSiteRequestPath(root, path string) string {
	if isCMSDomainRoot(root) {
		return mapCMSDomainRequestPath(path)
	}
	if isImagesDomainRoot(root) {
		prefix := cfg.GetImageStaticPrefix()
		if prefix == "" {
			prefix = "/images"
		}
		return mapLegacyPrefixRequestPath(path, prefix)
	}
	return path
}

func isImagesDomainRoot(root string) bool {
	root = filepath.Clean(root)
	if strings.EqualFold(filepath.Base(root), "images") {
		return true
	}
	imgRoot := cfg.GetStaticPathRoot("/images")
	return imgRoot != "" && filepath.Clean(imgRoot) == root
}

func mapLegacyPrefixRequestPath(path, prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return path
	}
	prefix = normalizeStaticURLPrefix(prefix)
	if path == prefix || strings.HasPrefix(path, prefix+"/") {
		path = strings.TrimPrefix(path, prefix)
	}
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func normalizeStaticURLPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimRight(prefix, "/")
}

func isCMSDomainRoot(root string) bool {
	return strings.EqualFold(filepath.Base(filepath.Clean(root)), "cms")
}

func mapCMSDomainRequestPath(path string) string {
	if path == "/cms" || strings.HasPrefix(path, "/cms/") {
		path = strings.TrimPrefix(path, "/cms")
	}
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func serveCMSDomainSPA(r *ghttp.Request, root, reqPath string) {
	rel := strings.TrimPrefix(reqPath, "/")
	if rel != "" && isCMSAssetRequest(rel) {
		applyRequestCORS(r)
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	indexPath := filepath.Join(root, "index.html")
	if !gfile.Exists(indexPath) {
		applyRequestCORS(r)
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	writeStaticFile(r, indexPath)
	r.ExitAll()
}

func buildDomainTLSConfig(ctx context.Context) *tls.Config {
	if len(domainCertEntries) == 0 {
		certPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsCertPath").String())
		keyPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsKeyPath").String())
		if certPath == "" || keyPath == "" {
			return nil
		}
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			g.Log().Errorf(ctx, "加载默认HTTPS证书失败 cert=%s key=%s err=%v", certPath, keyPath, err)
			return nil
		}
		return &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
		}
	}
	certs := make([]tls.Certificate, 0, len(domainCertEntries))
	for _, item := range domainCertEntries {
		certs = append(certs, item.cert)
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: certs,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if info == nil {
				return nil, fmt.Errorf("empty tls client hello")
			}
			name := strings.ToLower(strings.TrimSpace(info.ServerName))
			for i := range domainCertEntries {
				for _, domain := range domainCertEntries[i].domains {
					if strings.ToLower(domain) == name {
						return &domainCertEntries[i].cert, nil
					}
				}
			}
			if len(certs) > 0 {
				return &certs[0], nil
			}
			return nil, fmt.Errorf("no certificate configured for host %s", info.ServerName)
		},
	}
}

func bindCMSStaticFallback(ctx context.Context) {
	if len(cfg.GetStaticPathCfgs()) > 0 {
		return
	}
	root := strings.TrimSpace(cfg.GetServerRoot())
	if root == "" {
		return
	}
	indexPath := filepath.Join(filepath.Clean(root), "cms", "index.html")
	if !gfile.Exists(indexPath) {
		return
	}
	httpServer.BindHookHandler("/cms/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		serveCMSStaticFallback(r, root, indexPath)
	})
	httpServer.BindHookHandler("/cms", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		serveCMSStaticFallback(r, root, indexPath)
	})
	g.Log().Warning(ctx, "已启用 CMS SPA 路由回退 /cms -> index.html")
}

func serveCMSStaticFallback(r *ghttp.Request, root, indexPath string) {
	if r == nil {
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}
	reqPath := r.URL.Path
	if reqPath != "/cms" && !strings.HasPrefix(reqPath, "/cms/") {
		return
	}
	filePath, ok := buildDomainStaticFilePath(root, reqPath)
	if ok && gfile.Exists(filePath) && !gfile.IsDir(filePath) {
		return
	}
	writeStaticFile(r, indexPath)
	r.ExitAll()
}

func buildDomainStaticFilePath(root, reqPath string) (string, bool) {
	if reqPath == "" {
		reqPath = "/"
	}
	cleanRoot := filepath.Clean(root)
	target := filepath.Clean(filepath.Join(cleanRoot, filepath.FromSlash(reqPath)))
	if target != cleanRoot && !strings.HasPrefix(target, cleanRoot+string(filepath.Separator)) {
		return "", false
	}
	if gfile.IsDir(target) {
		target = filepath.Join(target, "index.html")
	}
	return target, true
}
