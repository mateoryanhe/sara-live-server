package domainsite

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/entity/sys"
)

// Init 将数据库中的站点映射加载到进程缓存。
func Init() {
	for _, row := range domainsitedao.ListMappings() {
		if row == nil || strings.TrimSpace(row.SiteKey) == "" {
			continue
		}
		cfg.RegisterRuntimeStaticSite(row.Prefix, row.Domain, row.Root)
	}
}

// SyncHTTPSRegistryFromDB 按 domain_site_mappings 重建 Host→目录注册表(HTTP 启动前调用)。
func SyncHTTPSRegistryFromDB(ctx context.Context) error {
	return syncDomainSiteRegistryFromDB(ctx)
}

func syncDomainSiteRegistryFromDB(ctx context.Context) error {
	defaultCert, defaultKey := resolveDefaultTLSFiles(ctx)
	byDomain := make(map[string]httpserver.DomainSiteRegistration)
	for _, row := range domainsitedao.ListMappings() {
		if row == nil {
			continue
		}
		domain, err := normalizeDomain(row.Domain)
		if err != nil {
			continue
		}
		root, err := normalizeRoot(row.Root)
		if err != nil {
			continue
		}
		certFile, keyFile := strings.TrimSpace(row.CertFile), strings.TrimSpace(row.KeyFile)
		if certFile == "" || keyFile == "" || !gfile.Exists(certFile) || !gfile.Exists(keyFile) {
			certFile, keyFile = defaultCert, defaultKey
		}
		byDomain[domain] = httpserver.DomainSiteRegistration{
			Domain:   domain,
			Root:     root,
			CertFile: certFile,
			KeyFile:  keyFile,
		}
	}
	if len(byDomain) == 0 {
		return nil
	}
	registrations := make([]httpserver.DomainSiteRegistration, 0, len(byDomain))
	for _, item := range byDomain {
		registrations = append(registrations, item)
	}
	_, err := httpserver.RefreshDomainSiteRegistry(ctx, registrations)
	return err
}

func defaultDomainSiteTLSFiles(ctx context.Context) (string, string) {
	if ctx == nil {
		ctx = context.Background()
	}
	return strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsCertPath").String()),
		strings.TrimSpace(g.Cfg().MustGet(ctx, "server.httpsKeyPath").String())
}

// resolveDefaultTLSFiles 读取 config.yaml 中 Cloudflare 源站证书路径;文件不存在则返回空。
func resolveDefaultTLSFiles(ctx context.Context) (certFile, keyFile string) {
	certFile, keyFile = defaultDomainSiteTLSFiles(ctx)
	if certFile == "" || keyFile == "" || !gfile.Exists(certFile) || !gfile.Exists(keyFile) {
		return "", ""
	}
	return certFile, keyFile
}

// SaveStaticSiteMapping 保存一个业务站点的唯一域名映射，并热刷新当前进程。
func SaveStaticSiteMapping(ctx context.Context, siteKey, prefix, domain, root string) (*entity.DomainSiteMapping, error) {
	return SaveStaticSiteMappingWithTLS(ctx, siteKey, prefix, domain, root, "", "")
}

// SaveStaticSiteMappingWithTLS 同 SaveStaticSiteMapping,可指定证书路径;空则保留已有或回落 config.yaml 默认源站证书。
func SaveStaticSiteMappingWithTLS(ctx context.Context, siteKey, prefix, domain, root, certFile, keyFile string) (*entity.DomainSiteMapping, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	siteKey = strings.TrimSpace(siteKey)
	prefix = normalizePrefix(prefix)
	domain, err := normalizeDomain(domain)
	if err != nil {
		return nil, err
	}
	root, err = normalizeRoot(root)
	if err != nil {
		return nil, err
	}
	if siteKey == "" || prefix == "" || prefix == "/" {
		return nil, fmt.Errorf("站点键、URL前缀、域名和部署目录不能为空")
	}

	existing := domainsitedao.GetMappingBySiteKey(siteKey)
	oldDomains := make([]string, 0, 1)
	if existing != nil && strings.TrimSpace(existing.Domain) != "" {
		oldDomains = append(oldDomains, existing.Domain)
	} else {
		oldDomains = cfg.SplitDomains(cfg.GetStaticSiteDomains(prefix))
	}

	previousSnapshot := httpserver.GetDomainSiteRegistrations()
	oldDomainSet := make(map[string]struct{}, len(oldDomains))
	for _, item := range oldDomains {
		item = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(item), "."))
		if item != "" {
			oldDomainSet[item] = struct{}{}
		}
	}
	certFile = strings.TrimSpace(certFile)
	keyFile = strings.TrimSpace(keyFile)
	if certFile == "" || keyFile == "" {
		if existing != nil {
			certFile = strings.TrimSpace(existing.CertFile)
			keyFile = strings.TrimSpace(existing.KeyFile)
		}
	}
	if certFile == "" || keyFile == "" {
		for _, item := range previousSnapshot {
			if _, ok := oldDomainSet[item.Domain]; !ok {
				continue
			}
			certFile = item.CertFile
			keyFile = item.KeyFile
			break
		}
	}
	if certFile == "" || keyFile == "" || !gfile.Exists(certFile) || !gfile.Exists(keyFile) {
		certFile, keyFile = resolveDefaultTLSFiles(ctx)
	}

	if _, err = httpserver.ReplaceDomainSiteRegistration(ctx, oldDomains, httpserver.DomainSiteRegistration{
		Domain:   domain,
		Root:     root,
		CertFile: certFile,
		KeyFile:  keyFile,
	}); err != nil {
		return nil, err
	}

	row := existing
	if row == nil {
		row = &entity.DomainSiteMapping{}
	}
	row.SiteKey = siteKey
	row.Prefix = prefix
	row.Domain = domain
	row.Root = root
	row.CertFile = certFile
	row.KeyFile = keyFile
	row.UpdatedAt = time.Now()
	if err = domainsitedao.SaveMapping(row); err != nil {
		if _, rollbackErr := httpserver.RefreshDomainSiteRegistry(ctx, previousSnapshot); rollbackErr != nil {
			g.Log().Errorf(ctx, "恢复域名注册表失败 siteKey=%s err=%v", siteKey, rollbackErr)
		}
		return nil, err
	}

	cfg.RegisterRuntimeStaticSite(prefix, domain, root)
	g.Log().Infof(ctx, "动态站点映射保存并刷新成功 siteKey=%s prefix=%s domain=%s root=%s cert=%s", siteKey, prefix, domain, root, certFile)
	return row, nil
}

func normalizePrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimRight(prefix, "/")
}

func normalizeDomain(domain string) (string, error) {
	if strings.Contains(domain, ",") {
		return "", fmt.Errorf("只允许配置一个域名,不能使用逗号")
	}
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if domain == "" {
		return "", fmt.Errorf("域名不能为空")
	}
	return domain, nil
}

func normalizeRoot(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "", fmt.Errorf("部署目录不能为空")
	}
	if strings.HasPrefix(root, "/") {
		return path.Clean(root), nil
	}
	if !filepath.IsAbs(root) {
		return "", fmt.Errorf("部署目录必须是绝对路径: %s", root)
	}
	return filepath.Clean(root), nil
}
