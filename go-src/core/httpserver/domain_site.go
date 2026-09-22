package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/cfg"
)

// DomainSiteRegistration 是一项可热刷新的域名静态目录映射。
// Domain 兼容逗号分隔的多个域名；请求匹配时统一按不带端口的小写 Host 精确匹配。
type DomainSiteRegistration struct {
	Domain   string
	Root     string
	CertFile string
	KeyFile  string
}

type domainSiteRegistrySnapshot struct {
	roots            map[string]string
	certificates     map[string]*tls.Certificate
	registrations    []DomainSiteRegistration
	firstCertificate *tls.Certificate
}

var (
	domainSiteRegistry       atomic.Pointer[domainSiteRegistrySnapshot]
	domainSiteRegistryMu     sync.Mutex
	domainStaticHookBindOnce sync.Once
	defaultTLSCertificate    *tls.Certificate
)

func setupDomainSites() {
	ctx := gctx.New()
	bindDomainStaticHook()
	active, err := RefreshDomainSiteRegistryFromConfig(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "初始化域名静态目录注册表失败 err=%v", err)
	}
	if tlsConfig := buildDomainTLSConfig(ctx); tlsConfig != nil {
		httpServer.SetTLSConfig(tlsConfig)
		g.Log().Warning(ctx, "已启用多域名HTTPS证书(SNI)")
	}
	if len(active) > 0 {
		g.Log().Warningf(ctx, "已启用可热刷新的域名静态目录注册表,共 %d 个域名", len(active))
	}
}

// RefreshDomainSiteRegistryFromConfig 使用当前配置与运行时注册项重建域名注册表。
func RefreshDomainSiteRegistryFromConfig(ctx context.Context) ([]DomainSiteRegistration, error) {
	sites := cfg.GetDomainSiteCfgs()
	registrations := make([]DomainSiteRegistration, 0, len(sites))
	for _, site := range sites {
		if site == nil {
			continue
		}
		registrations = append(registrations, DomainSiteRegistration{
			Domain:   site.Domain,
			Root:     site.Root,
			CertFile: site.CertFile,
			KeyFile:  site.KeyFile,
		})
	}
	return RefreshDomainSiteRegistry(ctx, registrations)
}

// RefreshDomainSiteRegistry 校验完整数组后原子替换注册表。
// 校验失败时保留旧快照，正在处理的请求不会读到半份配置。
func RefreshDomainSiteRegistry(ctx context.Context, registrations []DomainSiteRegistration) ([]DomainSiteRegistration, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	domainSiteRegistryMu.Lock()
	defer domainSiteRegistryMu.Unlock()
	return refreshDomainSiteRegistryLocked(ctx, registrations)
}

// ReplaceDomainSiteRegistration 以域名为键替换单个站点映射，并原子发布新快照。
// previousDomains 可传旧配置中的逗号分隔域名；未命中的其他站点保持不变。
func ReplaceDomainSiteRegistration(ctx context.Context, previousDomains []string, replacement DomainSiteRegistration) ([]DomainSiteRegistration, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	oldDomainSet := make(map[string]struct{})
	for _, raw := range previousDomains {
		for _, item := range cfg.SplitDomains(raw) {
			domain, err := normalizeRegisteredDomain(item)
			if err != nil {
				return nil, err
			}
			oldDomainSet[domain] = struct{}{}
		}
	}

	domainSiteRegistryMu.Lock()
	defer domainSiteRegistryMu.Unlock()
	current := domainSiteRegistry.Load()
	registrations := make([]DomainSiteRegistration, 0, len(oldDomainSet)+1)
	if current != nil {
		registrations = make([]DomainSiteRegistration, 0, len(current.registrations)+1)
		for _, item := range current.registrations {
			if _, replaced := oldDomainSet[item.Domain]; replaced {
				continue
			}
			registrations = append(registrations, item)
		}
	}
	registrations = append(registrations, replacement)
	return refreshDomainSiteRegistryLocked(ctx, registrations)
}

func refreshDomainSiteRegistryLocked(ctx context.Context, registrations []DomainSiteRegistration) ([]DomainSiteRegistration, error) {
	snapshot, err := buildDomainSiteRegistrySnapshot(ctx, registrations)
	if err != nil {
		return nil, err
	}
	domainSiteRegistry.Store(snapshot)
	g.Log().Infof(ctx, "域名静态目录注册表刷新成功,domains=%d", len(snapshot.registrations))
	return cloneDomainSiteRegistrations(snapshot.registrations), nil
}

// GetDomainSiteRegistrations 返回当前进程正在使用的域名注册表快照。
func GetDomainSiteRegistrations() []DomainSiteRegistration {
	snapshot := domainSiteRegistry.Load()
	if snapshot == nil {
		return []DomainSiteRegistration{}
	}
	return cloneDomainSiteRegistrations(snapshot.registrations)
}

func buildDomainSiteRegistrySnapshot(ctx context.Context, registrations []DomainSiteRegistration) (*domainSiteRegistrySnapshot, error) {
	snapshot := &domainSiteRegistrySnapshot{
		roots:         make(map[string]string),
		certificates:  make(map[string]*tls.Certificate),
		registrations: make([]DomainSiteRegistration, 0, len(registrations)),
	}
	for index, registration := range registrations {
		root := resolveDomainSiteRoot(ctx, registration.Root)
		if root == "" {
			return nil, fmt.Errorf("第 %d 项域名站点根目录不能为空", index+1)
		}
		domains := cfg.SplitDomains(registration.Domain)
		if len(domains) == 0 {
			return nil, fmt.Errorf("第 %d 项域名不能为空", index+1)
		}
		certFile := strings.TrimSpace(registration.CertFile)
		keyFile := strings.TrimSpace(registration.KeyFile)
		if (certFile == "") != (keyFile == "") {
			return nil, fmt.Errorf("域名 %s 的证书文件和私钥文件必须同时配置", registration.Domain)
		}
		var cert *tls.Certificate
		if certFile != "" {
			loaded, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("加载域名 %s 的证书失败: %w", registration.Domain, err)
			}
			cert = &loaded
			if snapshot.firstCertificate == nil {
				snapshot.firstCertificate = cert
			}
		}
		for _, rawDomain := range domains {
			domain, err := normalizeRegisteredDomain(rawDomain)
			if err != nil {
				return nil, fmt.Errorf("第 %d 项域名无效: %w", index+1, err)
			}
			if _, exists := snapshot.roots[domain]; exists {
				return nil, fmt.Errorf("域名重复注册: %s", domain)
			}
			snapshot.roots[domain] = root
			if cert != nil {
				snapshot.certificates[domain] = cert
			}
			snapshot.registrations = append(snapshot.registrations, DomainSiteRegistration{
				Domain:   domain,
				Root:     root,
				CertFile: certFile,
				KeyFile:  keyFile,
			})
		}
	}
	return snapshot, nil
}

func normalizeRegisteredDomain(domain string) (string, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	domain = strings.TrimSuffix(domain, ".")
	if domain == "" {
		return "", fmt.Errorf("域名不能为空")
	}
	if strings.Contains(domain, "://") || strings.ContainsAny(domain, `/\\?#@`) {
		return "", fmt.Errorf("只允许填写 Host,不能包含协议、路径或查询参数: %s", domain)
	}
	if strings.ContainsAny(domain, " \t\r\n") {
		return "", fmt.Errorf("域名不能包含空白字符: %s", domain)
	}
	if strings.Contains(domain, ":") {
		return "", fmt.Errorf("域名不能包含端口: %s", domain)
	}
	if net.ParseIP(domain) != nil || domain == "localhost" {
		return domain, nil
	}
	if len(domain) > 253 {
		return "", fmt.Errorf("域名长度超过253字符: %s", domain)
	}
	for _, label := range strings.Split(domain, ".") {
		if len(label) == 0 || len(label) > 63 {
			return "", fmt.Errorf("域名标签长度无效: %s", domain)
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return "", fmt.Errorf("域名标签不能以连字符开头或结尾: %s", domain)
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return "", fmt.Errorf("域名包含非法字符: %s", domain)
			}
		}
	}
	return domain, nil
}

func cloneDomainSiteRegistrations(list []DomainSiteRegistration) []DomainSiteRegistration {
	if len(list) == 0 {
		return []DomainSiteRegistration{}
	}
	cloned := make([]DomainSiteRegistration, len(list))
	copy(cloned, list)
	return cloned
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

func bindDomainStaticHook() {
	domainStaticHookBindOnce.Do(func() {
		httpServer.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
			if r == nil {
				return
			}
			snapshot := domainSiteRegistry.Load()
			if snapshot == nil {
				return
			}
			host := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(r.GetHost()), "."))
			root := snapshot.roots[host]
			if root == "" {
				return
			}
			serveDomainStatic(r, root)
		})
	})
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
	certPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsCertPath").String())
	keyPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsKeyPath").String())
	if certPath != "" && keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			g.Log().Errorf(ctx, "加载默认HTTPS证书失败 cert=%s key=%s err=%v", certPath, keyPath, err)
		} else {
			defaultTLSCertificate = &cert
		}
	}
	snapshot := domainSiteRegistry.Load()
	firstCertificate := defaultTLSCertificate
	if firstCertificate == nil && snapshot != nil {
		firstCertificate = snapshot.firstCertificate
	}
	if firstCertificate == nil {
		return nil
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{*firstCertificate},
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if info == nil {
				return nil, fmt.Errorf("empty tls client hello")
			}
			name := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(info.ServerName), "."))
			current := domainSiteRegistry.Load()
			if current != nil {
				if cert := current.certificates[name]; cert != nil {
					return cert, nil
				}
			}
			if defaultTLSCertificate != nil {
				return defaultTLSCertificate, nil
			}
			if current != nil && current.firstCertificate != nil {
				return current.firstCertificate, nil
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
