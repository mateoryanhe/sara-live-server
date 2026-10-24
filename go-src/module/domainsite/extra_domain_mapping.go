package domainsite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/domainsitedto"
	sysentity "xr-game-server/entity/sys"
)

const extraDomainSiteKeyPrefix = "extra-domain-"

func extraDomainSiteKey(domain string) string {
	sum := sha256.Sum256([]byte(domain))
	return extraDomainSiteKeyPrefix + hex.EncodeToString(sum[:8])
}

func extraDomainURLPrefix(domain string) string {
	sum := sha256.Sum256([]byte(domain))
	return "/extra-domain/" + hex.EncodeToString(sum[:8])
}

func isExtraDomainSiteKey(siteKey string) bool {
	return strings.HasPrefix(strings.TrimSpace(siteKey), extraDomainSiteKeyPrefix)
}

func ListExtraDomainMappings(_ context.Context, _ *domainsitedto.ListExtraDomainMappingsReq) (*domainsitedto.ListExtraDomainMappingsRes, error) {
	items := make([]*domainsitedto.ExtraDomainMappingItem, 0)
	for _, row := range domainsitedao.ListMappings() {
		if row == nil || !isExtraDomainSiteKey(row.SiteKey) {
			continue
		}
		items = append(items, extraDomainMappingItem(row))
	}
	return &domainsitedto.ListExtraDomainMappingsRes{Items: items}, nil
}

func SaveExtraDomainMapping(ctx context.Context, req *domainsitedto.SaveExtraDomainMappingReq) (*domainsitedto.SaveExtraDomainMappingRes, error) {
	if req == nil {
		return nil, fmt.Errorf("request is empty")
	}
	domain, err := normalizeDomain(req.Domain)
	if err != nil {
		return nil, err
	}
	root, err := normalizeRoot(req.Root)
	if err != nil {
		return nil, err
	}
	siteKey := extraDomainSiteKey(domain)
	if domainUsedByOtherMapping(domain, siteKey) {
		return nil, fmt.Errorf("域名已被其他映射使用: %s", domain)
	}
	row, err := SaveStaticSiteMapping(ctx, siteKey, extraDomainURLPrefix(domain), domain, root)
	if err != nil {
		return nil, err
	}
	return &domainsitedto.SaveExtraDomainMappingRes{
		Success: true,
		Item:    extraDomainMappingItem(row),
	}, nil
}

// SaveExtraDomainMappingFromRequest multipart: domain + root + cert + key(证书/私钥必传)。
func SaveExtraDomainMappingFromRequest(r *ghttp.Request) (*domainsitedto.SaveExtraDomainMappingRes, error) {
	if r == nil {
		return nil, fmt.Errorf("request is empty")
	}
	domain, err := normalizeDomain(r.Get("domain").String())
	if err != nil {
		return nil, err
	}
	root, err := normalizeRoot(r.Get("root").String())
	if err != nil {
		return nil, err
	}
	siteKey := extraDomainSiteKey(domain)
	if domainUsedByOtherMapping(domain, siteKey) {
		return nil, fmt.Errorf("域名已被其他映射使用: %s", domain)
	}
	certPart := r.GetUploadFile("cert")
	keyPart := r.GetUploadFile("key")
	if certPart == nil {
		certPart = r.GetUploadFile("certFile")
	}
	if keyPart == nil {
		keyPart = r.GetUploadFile("keyFile")
	}
	certPath, keyPath, err := storeExtraDomainTLSFiles(siteKey, certPart, keyPart)
	if err != nil {
		return nil, err
	}
	row, err := SaveStaticSiteMappingWithTLS(
		r.Context(),
		siteKey,
		extraDomainURLPrefix(domain),
		domain,
		root,
		certPath,
		keyPath,
	)
	if err != nil {
		removeExtraDomainTLSDir(siteKey)
		return nil, err
	}
	return &domainsitedto.SaveExtraDomainMappingRes{
		Success: true,
		Item:    extraDomainMappingItem(row),
	}, nil
}

func DeleteExtraDomainMapping(ctx context.Context, req *domainsitedto.DeleteExtraDomainMappingReq) (*domainsitedto.DeleteExtraDomainMappingRes, error) {
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
	if !isExtraDomainSiteKey(row.SiteKey) {
		return nil, fmt.Errorf("不是额外域名映射记录")
	}
	previousSnapshot := httpserver.GetDomainSiteRegistrations()
	if err = removeDomainFromRegistry(ctx, row.Domain); err != nil {
		return nil, err
	}
	if err = domainsitedao.DeleteMapping(id); err != nil {
		if _, rollbackErr := httpserver.RefreshDomainSiteRegistry(ctx, previousSnapshot); rollbackErr != nil {
			g.Log().Errorf(ctx, "回滚额外域名映射删除失败 id=%d err=%v", id, rollbackErr)
		}
		return nil, err
	}
	removeExtraDomainTLSDir(row.SiteKey)
	g.Log().Infof(ctx, "额外域名映射已删除 domain=%s id=%d", row.Domain, id)
	return &domainsitedto.DeleteExtraDomainMappingRes{Success: true}, nil
}

func extraDomainMappingItem(row *sysentity.DomainSiteMapping) *domainsitedto.ExtraDomainMappingItem {
	if row == nil {
		return nil
	}
	cert := strings.TrimSpace(row.CertFile)
	key := strings.TrimSpace(row.KeyFile)
	hasTLS := cert != "" && key != "" && gfile.Exists(cert) && gfile.Exists(key)
	return &domainsitedto.ExtraDomainMappingItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Domain:    row.Domain,
		Root:      row.Root,
		HasTLS:    hasTLS,
		CertFile:  cert,
		KeyFile:   key,
		UpdatedAt: formatDomainSiteTime(row.UpdatedAt),
	}
}

func domainUsedByOtherMapping(domain, exceptSiteKey string) bool {
	for _, row := range domainsitedao.ListMappings() {
		if row == nil || row.SiteKey == exceptSiteKey {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(row.Domain), domain) {
			return true
		}
	}
	return false
}

func removeDomainFromRegistry(ctx context.Context, domain string) error {
	target, err := normalizeDomain(domain)
	if err != nil {
		return err
	}
	current := httpserver.GetDomainSiteRegistrations()
	next := make([]httpserver.DomainSiteRegistration, 0, len(current))
	for _, item := range current {
		d, normErr := normalizeDomain(item.Domain)
		if normErr != nil {
			next = append(next, item)
			continue
		}
		if d == target {
			continue
		}
		next = append(next, item)
	}
	_, err = httpserver.RefreshDomainSiteRegistry(ctx, next)
	return err
}
