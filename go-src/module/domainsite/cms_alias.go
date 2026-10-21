package domainsite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/domainsitedto"
	sysentity "xr-game-server/entity/sys"
)

const cmsAliasSiteKeyPrefix = "cms-alias-"

func cmsAliasSiteKey(domain string) string {
	sum := sha256.Sum256([]byte(domain))
	return cmsAliasSiteKeyPrefix + hex.EncodeToString(sum[:8])
}

func isCMSAliasSiteKey(siteKey string) bool {
	return strings.HasPrefix(strings.TrimSpace(siteKey), cmsAliasSiteKeyPrefix)
}

func cmsDeployRoot() (string, error) {
	if row := domainsitedao.GetMappingBySiteKey(domainsitedto.CMSSiteKey); row != nil && strings.TrimSpace(row.Root) != "" {
		return normalizeRoot(row.Root)
	}
	root := strings.TrimSpace(cfg.GetStaticPathRoot(domainsitedto.CMSURLPrefix))
	if root == "" {
		root = domainsitedto.DefaultCMSRoot
	}
	return normalizeRoot(root)
}

func cmsPrimaryDomain() string {
	if row := domainsitedao.GetMappingBySiteKey(domainsitedto.CMSSiteKey); row != nil {
		return strings.TrimSpace(row.Domain)
	}
	return strings.TrimSpace(cfg.GetStaticSiteDomain(domainsitedto.CMSURLPrefix))
}

func certFilesForDomain(domain string) (certFile, keyFile string) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	for _, item := range httpserver.GetDomainSiteRegistrations() {
		for _, d := range cfg.SplitDomains(item.Domain) {
			if strings.ToLower(strings.TrimSpace(d)) == domain {
				return item.CertFile, item.KeyFile
			}
		}
	}
	return "", ""
}

func removeDomainsFromRegistry(ctx context.Context, domains ...string) error {
	removeSet := make(map[string]struct{}, len(domains))
	for _, raw := range domains {
		d, err := normalizeDomain(raw)
		if err != nil {
			return err
		}
		removeSet[d] = struct{}{}
	}
	current := httpserver.GetDomainSiteRegistrations()
	next := make([]httpserver.DomainSiteRegistration, 0, len(current))
	for _, item := range current {
		d, err := normalizeDomain(item.Domain)
		if err != nil {
			continue
		}
		if _, drop := removeSet[d]; drop {
			continue
		}
		next = append(next, item)
	}
	_, err := httpserver.RefreshDomainSiteRegistry(ctx, next)
	return err
}

func listCMSAliasRows() []*sysentity.DomainSiteMapping {
	rows := domainsitedao.ListMappings()
	out := make([]*sysentity.DomainSiteMapping, 0)
	for _, row := range rows {
		if row != nil && isCMSAliasSiteKey(row.SiteKey) {
			out = append(out, row)
		}
	}
	return out
}

func aliasItem(row *sysentity.DomainSiteMapping) *domainsitedto.CMSDomainAliasItem {
	if row == nil {
		return nil
	}
	return &domainsitedto.CMSDomainAliasItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Domain:    row.Domain,
		Root:      row.Root,
		UpdatedAt: formatDomainSiteTime(row.UpdatedAt),
	}
}

func domainTakenByCMS(domain string) bool {
	primary := strings.ToLower(strings.TrimSpace(cmsPrimaryDomain()))
	if primary != "" && primary == domain {
		return true
	}
	for _, row := range listCMSAliasRows() {
		if row != nil && strings.EqualFold(strings.TrimSpace(row.Domain), domain) {
			return true
		}
	}
	return false
}

// ListCMSDomainAlias CMS 备用 CMS 访问域名(与主映射共用同一静态目录).
func ListCMSDomainAlias(_ context.Context, _ *domainsitedto.ListCMSDomainAliasReq) (*domainsitedto.ListCMSDomainAliasRes, error) {
	root, err := cmsDeployRoot()
	if err != nil {
		return nil, err
	}
	items := make([]*domainsitedto.CMSDomainAliasItem, 0)
	for _, row := range listCMSAliasRows() {
		items = append(items, aliasItem(row))
	}
	return &domainsitedto.ListCMSDomainAliasRes{
		PrimaryDomain: cmsPrimaryDomain(),
		Root:          root,
		Items:         items,
	}, nil
}

// AddCMSDomainAlias 追加一个 Host,热刷新域名注册表.
func AddCMSDomainAlias(ctx context.Context, req *domainsitedto.AddCMSDomainAliasReq) (*domainsitedto.AddCMSDomainAliasRes, error) {
	if req == nil {
		return nil, fmt.Errorf("request is empty")
	}
	domain, err := normalizeDomain(req.Domain)
	if err != nil {
		return nil, err
	}
	if domainTakenByCMS(domain) {
		return nil, fmt.Errorf("域名已存在: %s", domain)
	}
	root, err := cmsDeployRoot()
	if err != nil {
		return nil, err
	}
	siteKey := cmsAliasSiteKey(domain)
	if existing := domainsitedao.GetMappingBySiteKey(siteKey); existing != nil && existing.ID != 0 {
		return nil, fmt.Errorf("域名已存在: %s", domain)
	}

	previousSnapshot := httpserver.GetDomainSiteRegistrations()
	primary := strings.ToLower(strings.TrimSpace(cmsPrimaryDomain()))
	certFile, keyFile := certFilesForDomain(primary)
	if certFile == "" && primary != "" {
		certFile, keyFile = certFilesForDomain(domain)
	}
	if _, err = httpserver.ReplaceDomainSiteRegistration(ctx, nil, httpserver.DomainSiteRegistration{
		Domain:   domain,
		Root:     root,
		CertFile: certFile,
		KeyFile:  keyFile,
	}); err != nil {
		return nil, err
	}

	row := &sysentity.DomainSiteMapping{
		SiteKey: siteKey,
		Prefix:  domainsitedto.CMSURLPrefix,
		Domain:  domain,
		Root:    root,
	}
	row.UpdatedAt = time.Now()
	if err = domainsitedao.SaveMapping(row); err != nil {
		if _, rollbackErr := httpserver.RefreshDomainSiteRegistry(ctx, previousSnapshot); rollbackErr != nil {
			g.Log().Errorf(ctx, "回滚CMS备用域名注册表失败 domain=%s err=%v", domain, rollbackErr)
		}
		return nil, err
	}
	cfg.RegisterRuntimeStaticSite(domainsitedto.CMSURLPrefix, domain, root)
	g.Log().Infof(ctx, "CMS备用域名已添加 domain=%s root=%s", domain, root)
	return &domainsitedto.AddCMSDomainAliasRes{Success: true, Item: aliasItem(row)}, nil
}

// DeleteCMSDomainAlias 删除备用域名并从注册表移除.
func DeleteCMSDomainAlias(ctx context.Context, req *domainsitedto.DeleteCMSDomainAliasReq) (*domainsitedto.DeleteCMSDomainAliasRes, error) {
	if req == nil {
		return nil, fmt.Errorf("request is empty")
	}
	id, err := strconv.ParseUint(strings.TrimSpace(req.ID), 10, 64)
	if err != nil || id == 0 {
		return nil, fmt.Errorf("无效的记录ID")
	}
	var row sysentity.DomainSiteMapping
	if err = g.DB().Model(string(sysentity.TbDomainSiteMapping)).WherePri(id).Scan(&row); err != nil || row.ID == 0 {
		return nil, fmt.Errorf("记录不存在")
	}
	if !isCMSAliasSiteKey(row.SiteKey) {
		return nil, fmt.Errorf("不是CMS备用域名记录")
	}

	previousSnapshot := httpserver.GetDomainSiteRegistrations()
	if err = removeDomainsFromRegistry(ctx, row.Domain); err != nil {
		return nil, err
	}
	if err = domainsitedao.DeleteMapping(id); err != nil {
		if _, rollbackErr := httpserver.RefreshDomainSiteRegistry(ctx, previousSnapshot); rollbackErr != nil {
			g.Log().Errorf(ctx, "回滚CMS备用域名删除失败 id=%d err=%v", id, rollbackErr)
		}
		return nil, err
	}
	g.Log().Infof(ctx, "CMS备用域名已删除 domain=%s id=%d", row.Domain, id)
	return &domainsitedto.DeleteCMSDomainAliasRes{Success: true}, nil
}

// syncCMSAliasRoots 主 CMS 目录变更时,同步所有备用域名记录的 root 与注册表.
func syncCMSAliasRoots(ctx context.Context, root string) {
	root = strings.TrimSpace(root)
	if root == "" {
		return
	}
	rows := listCMSAliasRows()
	if len(rows) == 0 {
		return
	}
	snapshot := httpserver.GetDomainSiteRegistrations()
	updated := make([]httpserver.DomainSiteRegistration, 0, len(snapshot))
	for _, item := range snapshot {
		d, normErr := normalizeDomain(item.Domain)
		if normErr != nil {
			updated = append(updated, item)
			continue
		}
		isAlias := false
		for _, row := range rows {
			if row != nil && strings.EqualFold(strings.TrimSpace(row.Domain), d) {
				isAlias = true
				break
			}
		}
		if isAlias {
			item.Root = root
		}
		updated = append(updated, item)
	}
	if _, err := httpserver.RefreshDomainSiteRegistry(ctx, updated); err != nil {
		g.Log().Errorf(ctx, "同步CMS备用域名目录失败 root=%s err=%v", root, err)
		return
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.Root = root
		row.UpdatedAt = time.Now()
		if err := domainsitedao.SaveMapping(row); err != nil {
			g.Log().Errorf(ctx, "更新CMS备用域名目录失败 id=%d err=%v", row.ID, err)
		}
		cfg.RegisterRuntimeStaticSite(domainsitedto.CMSURLPrefix, row.Domain, root)
	}
}
