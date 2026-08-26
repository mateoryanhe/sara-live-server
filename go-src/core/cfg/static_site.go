package cfg

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// StaticSiteCfg 域名与物理目录统一配置(对外 HTTP 走 domain,对内路径解析走 prefix)
type StaticSiteCfg struct {
	Domain     string `json:"domain"     dc:"域名,多个用逗号分隔;本地开发可省略"`
	Prefix     string `json:"prefix"     dc:"URL前缀,如 /images;省略时从 domain 首段推导"`
	Path       string `json:"path"       dc:"物理目录绝对路径"`
	Root       string `json:"root"       dc:"物理目录,兼容旧字段"`
	TTLMinutes int    `json:"ttlMinutes" dc:"CMS文件导出过期清理(分钟),仅 cms-export 使用"`
	CertFile   string `json:"certFile"   dc:"HTTPS证书文件路径"`
	KeyFile    string `json:"keyFile"    dc:"HTTPS私钥文件路径"`
}

func staticSitePhysicalDir(item *StaticSiteCfg) string {
	if item == nil {
		return ""
	}
	if p := strings.TrimSpace(item.Path); p != "" {
		return p
	}
	return strings.TrimSpace(item.Root)
}

func initStaticSiteCfg() {
	ctx := gctx.New()
	sites := loadStaticSiteCfgs(ctx)
	staticPathCfgs = buildStaticPathCfgsFromSites(sites)
	domainSiteCfgs = buildDomainSiteCfgsFromSites(sites)
	initImageStaticCfg()
}

func loadStaticSiteCfgs(ctx context.Context) []*StaticSiteCfg {
	list := make([]*StaticSiteCfg, 0)
	_ = g.Cfg().MustGet(ctx, "server.staticSites").Scan(&list)
	if len(list) > 0 {
		return normalizeStaticSiteCfgs(list)
	}
	return mergeLegacyStaticSiteCfgs(ctx)
}

func mergeLegacyStaticSiteCfgs(ctx context.Context) []*StaticSiteCfg {
	legacyPaths := make([]*StaticPathCfg, 0)
	_ = g.Cfg().MustGet(ctx, "server.staticPaths").Scan(&legacyPaths)
	legacyDomains := make([]*DomainSiteCfg, 0)
	_ = g.Cfg().MustGet(ctx, "server.domainSites").Scan(&legacyDomains)

	byPath := make(map[string]*StaticSiteCfg)
	order := make([]string, 0)

	appendSite := func(site *StaticSiteCfg) {
		if site == nil {
			return
		}
		dir := staticSitePhysicalDir(site)
		if dir == "" {
			return
		}
		site.Path = dir
		site.Root = dir
		if _, ok := byPath[dir]; !ok {
			order = append(order, dir)
		}
		existing := byPath[dir]
		if existing == nil {
			byPath[dir] = site
			return
		}
		if existing.Domain == "" {
			existing.Domain = site.Domain
		}
		if existing.Prefix == "" {
			existing.Prefix = site.Prefix
		}
		if existing.TTLMinutes == 0 {
			existing.TTLMinutes = site.TTLMinutes
		}
		if existing.CertFile == "" {
			existing.CertFile = site.CertFile
		}
		if existing.KeyFile == "" {
			existing.KeyFile = site.KeyFile
		}
	}

	for _, item := range legacyPaths {
		if item == nil {
			continue
		}
		appendSite(&StaticSiteCfg{
			Prefix:     item.Prefix,
			Path:       staticPathPhysicalDir(item),
			TTLMinutes: item.TTLMinutes,
		})
	}
	for _, item := range legacyDomains {
		if item == nil {
			continue
		}
		appendSite(&StaticSiteCfg{
			Domain:   item.Domain,
			Path:     item.Root,
			CertFile: item.CertFile,
			KeyFile:  item.KeyFile,
		})
	}

	ret := make([]*StaticSiteCfg, 0, len(order))
	for _, dir := range order {
		if site := byPath[dir]; site != nil {
			ret = append(ret, site)
		}
	}
	return normalizeStaticSiteCfgs(ret)
}

func normalizeStaticSiteCfgs(list []*StaticSiteCfg) []*StaticSiteCfg {
	ret := make([]*StaticSiteCfg, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		item.Domain = strings.TrimSpace(item.Domain)
		item.Prefix = normalizeURLPrefix(strings.TrimSpace(item.Prefix))
		item.CertFile = strings.TrimSpace(item.CertFile)
		item.KeyFile = strings.TrimSpace(item.KeyFile)
		dir := staticSitePhysicalDir(item)
		if dir == "" {
			continue
		}
		item.Path = dir
		item.Root = dir
		if item.Prefix == "" || item.Prefix == "/" {
			if item.Domain != "" {
				item.Prefix = derivePrefixFromDomain(item.Domain)
			}
		}
		if item.Domain == "" && (item.Prefix == "" || item.Prefix == "/") {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func buildStaticPathCfgsFromSites(sites []*StaticSiteCfg) []*StaticPathCfg {
	list := make([]*StaticPathCfg, 0, len(sites))
	for _, site := range sites {
		if site == nil || site.Prefix == "" || site.Prefix == "/" {
			continue
		}
		list = append(list, &StaticPathCfg{
			Prefix:     site.Prefix,
			Path:       site.Path,
			Root:       site.Path,
			TTLMinutes: site.TTLMinutes,
		})
	}
	return normalizeStaticPathCfgs(list)
}

func buildDomainSiteCfgsFromSites(sites []*StaticSiteCfg) []*DomainSiteCfg {
	list := make([]*DomainSiteCfg, 0, len(sites))
	for _, site := range sites {
		if site == nil || site.Domain == "" {
			continue
		}
		list = append(list, &DomainSiteCfg{
			Domain:   site.Domain,
			Root:     site.Path,
			CertFile: site.CertFile,
			KeyFile:  site.KeyFile,
		})
	}
	return normalizeDomainSiteCfgs(list)
}

func derivePrefixFromDomain(domain string) string {
	domain = strings.TrimSpace(strings.Split(domain, ",")[0])
	if domain == "" {
		return ""
	}
	if idx := strings.Index(domain, ":"); idx > 0 {
		domain = domain[:idx]
	}
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}
	label := strings.ToLower(strings.TrimSpace(parts[0]))
	switch label {
	case "www":
		return ""
	}
	return normalizeURLPrefix(label)
}
