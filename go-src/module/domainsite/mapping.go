package domainsite

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/entity/sys"
)

// Init 将数据库中的站点映射加载到进程缓存。每个业务站点固定对应一条域名记录。
func Init() {
	for _, row := range domainsitedao.ListMappings() {
		if row == nil || strings.TrimSpace(row.SiteKey) == "" {
			continue
		}
		cfg.RegisterRuntimeStaticSite(row.Prefix, row.Domain, row.Root)
	}
}

// SaveStaticSiteMapping 保存一个业务站点的唯一域名映射，并热刷新当前进程。
func SaveStaticSiteMapping(ctx context.Context, siteKey, prefix, domain, root string) (*entity.DomainSiteMapping, error) {
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
	var certFile, keyFile string
	for _, item := range previousSnapshot {
		if _, ok := oldDomainSet[item.Domain]; !ok {
			continue
		}
		certFile = item.CertFile
		keyFile = item.KeyFile
		break
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
	row.UpdatedAt = time.Now()
	if err = domainsitedao.SaveMapping(row); err != nil {
		if _, rollbackErr := httpserver.RefreshDomainSiteRegistry(ctx, previousSnapshot); rollbackErr != nil {
			g.Log().Errorf(ctx, "恢复域名注册表失败 siteKey=%s err=%v", siteKey, rollbackErr)
		}
		return nil, err
	}

	cfg.RegisterRuntimeStaticSite(prefix, domain, root)
	g.Log().Infof(ctx, "动态站点映射保存并刷新成功 siteKey=%s prefix=%s domain=%s root=%s", siteKey, prefix, domain, root)
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
