package domainsite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/domainsitedto"
	"xr-game-server/entity/sys"
)

const hostTLSRoot = "/home/ec2-user/ssl/host-tls"

func hostTLSSiteKey(domain string) string {
	sum := sha256.Sum256([]byte(domain))
	return "host-tls-" + hex.EncodeToString(sum[:8])
}

func hostTLSDir(siteKey string) string {
	return filepath.Join(hostTLSRoot, strings.TrimSpace(siteKey))
}

func removeHostTLSDir(siteKey string) {
	dir := hostTLSDir(siteKey)
	if dir == "" || !strings.HasPrefix(filepath.Clean(dir), filepath.Clean(hostTLSRoot)) {
		return
	}
	_ = os.RemoveAll(dir)
}

// SyncHostTLSRegistryFromDB 启动时从 host_tls_mappings 重建纯 Host 证书表.
func SyncHostTLSRegistryFromDB(ctx context.Context) error {
	rows := domainsitedao.ListHostTlsMappings()
	regs := make([]httpserver.HostTLSRegistration, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		domain, err := normalizeHostTLSDomain(row.Domain)
		if err != nil {
			continue
		}
		cert := strings.TrimSpace(row.CertFile)
		key := strings.TrimSpace(row.KeyFile)
		if cert == "" || key == "" || !gfile.Exists(cert) || !gfile.Exists(key) {
			g.Log().Warningf(ctx, "跳过无效 Host 证书映射 domain=%s", domain)
			continue
		}
		regs = append(regs, httpserver.HostTLSRegistration{
			Domain:   domain,
			CertFile: cert,
			KeyFile:  key,
		})
	}
	return httpserver.RefreshHostTLSRegistry(ctx, regs)
}

func reloadHostTLSRegistryFromDB(ctx context.Context) error {
	return SyncHostTLSRegistryFromDB(ctx)
}

func ListHostTlsMappings(_ context.Context, _ *domainsitedto.ListHostTlsMappingsReq) (*domainsitedto.ListHostTlsMappingsRes, error) {
	items := make([]*domainsitedto.HostTlsMappingItem, 0)
	for _, row := range domainsitedao.ListHostTlsMappings() {
		items = append(items, hostTlsMappingItem(row))
	}
	return &domainsitedto.ListHostTlsMappingsRes{Items: items}, nil
}

// SaveHostTlsMappingFromRequest multipart: domain + cert + key(必传);无静态目录.
func SaveHostTlsMappingFromRequest(r *ghttp.Request) (*domainsitedto.SaveHostTlsMappingRes, error) {
	if r == nil {
		return nil, fmt.Errorf("request is empty")
	}
	ctx := r.Context()
	domain, err := normalizeHostTLSDomain(r.Get("domain").String())
	if err != nil {
		return nil, err
	}
	if domainUsedByStaticMapping(domain) {
		return nil, fmt.Errorf("域名已被静态目录映射占用: %s", domain)
	}
	certPart := r.GetUploadFile("cert")
	keyPart := r.GetUploadFile("key")
	if certPart == nil {
		certPart = r.GetUploadFile("certFile")
	}
	if keyPart == nil {
		keyPart = r.GetUploadFile("keyFile")
	}
	siteKey := hostTLSSiteKey(domain)
	certPath, keyPath, err := storeTLSFilesToDir(hostTLSDir(siteKey), certPart, keyPart)
	if err != nil {
		return nil, err
	}

	existing := domainsitedao.GetHostTlsMappingByDomain(domain)
	row := existing
	if row == nil {
		row = &entity.HostTlsMapping{}
	}
	row.Domain = domain
	row.CertFile = certPath
	row.KeyFile = keyPath
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err = domainsitedao.SaveHostTlsMapping(row); err != nil {
		removeHostTLSDir(siteKey)
		return nil, err
	}
	if err = reloadHostTLSRegistryFromDB(ctx); err != nil {
		g.Log().Errorf(ctx, "刷新纯Host证书注册表失败 domain=%s err=%v", domain, err)
		return nil, err
	}
	g.Log().Infof(ctx, "纯Host证书映射已保存 domain=%s cert=%s", domain, certPath)
	return &domainsitedto.SaveHostTlsMappingRes{
		Success: true,
		Item:    hostTlsMappingItem(row),
	}, nil
}

func DeleteHostTlsMapping(ctx context.Context, req *domainsitedto.DeleteHostTlsMappingReq) (*domainsitedto.DeleteHostTlsMappingRes, error) {
	if req == nil {
		return nil, fmt.Errorf("request is empty")
	}
	id, err := strconv.ParseUint(strings.TrimSpace(req.ID), 10, 64)
	if err != nil || id == 0 {
		return nil, fmt.Errorf("无效的记录ID")
	}
	row := domainsitedao.GetHostTlsMappingByID(id)
	if row == nil {
		return nil, fmt.Errorf("记录不存在")
	}
	domain := strings.TrimSpace(row.Domain)
	siteKey := hostTLSSiteKey(domain)
	if err = domainsitedao.DeleteHostTlsMapping(id); err != nil {
		return nil, err
	}
	removeHostTLSDir(siteKey)
	if err = reloadHostTLSRegistryFromDB(ctx); err != nil {
		g.Log().Errorf(ctx, "删除后刷新纯Host证书失败 domain=%s err=%v", domain, err)
		return nil, err
	}
	g.Log().Infof(ctx, "纯Host证书映射已删除 domain=%s id=%d", domain, id)
	return &domainsitedto.DeleteHostTlsMappingRes{Success: true}, nil
}

func hostTlsMappingItem(row *entity.HostTlsMapping) *domainsitedto.HostTlsMappingItem {
	if row == nil {
		return nil
	}
	cert := strings.TrimSpace(row.CertFile)
	key := strings.TrimSpace(row.KeyFile)
	hasTLS := cert != "" && key != "" && gfile.Exists(cert) && gfile.Exists(key)
	return &domainsitedto.HostTlsMappingItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Domain:    row.Domain,
		HasTLS:    hasTLS,
		CertFile:  cert,
		KeyFile:   key,
		UpdatedAt: formatDomainSiteTime(row.UpdatedAt),
	}
}

func normalizeHostTLSDomain(domain string) (string, error) {
	return httpserver.NormalizeHostTLSName(domain)
}

func domainUsedByStaticMapping(domain string) bool {
	for _, row := range domainsitedao.ListMappings() {
		if row == nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(row.Domain), domain) {
			return true
		}
	}
	return false
}

func domainUsedByHostTLS(domain string, exceptID uint64) bool {
	for _, row := range domainsitedao.ListHostTlsMappings() {
		if row == nil || row.ID == exceptID {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(row.Domain), domain) {
			return true
		}
	}
	return false
}
