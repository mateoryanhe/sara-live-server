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
		root := gfile.RealPath(site.Root)
		if root == "" {
			g.Log().Warningf(ctx, "域名站点根目录不存在,已跳过 domain=%s root=%s", site.Domain, site.Root)
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
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}
	filePath, ok := buildDomainStaticFilePath(root, r.URL.Path)
	if !ok || !gfile.Exists(filePath) {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	if r.Method == http.MethodHead {
		r.Response.WriteHeader(http.StatusOK)
	} else {
		r.Response.ServeFile(filePath)
	}
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
