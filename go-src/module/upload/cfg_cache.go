package upload

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/uploadresourcecfgdao"
	"xr-game-server/entity"
)

const (
	defaultResourceDomain = "http://127.0.0.1"
	defaultAvatarURL      = "https://img.yonogames.com/headimg/man/147.png"
)

type resourceCfgSnapshot struct {
	ResourceDomain   string
	DefaultAvatarUrl string
}

var (
	resourceCfgCache atomic.Value // *resourceCfgSnapshot
	emptyResourceCfg = &resourceCfgSnapshot{
		ResourceDomain:   defaultResourceDomain,
		DefaultAvatarUrl: defaultAvatarURL,
	}
)

func reloadResourceCfgMemory() {
	resourceCfgCache.Store(toResourceCfgSnapshot(uploadresourcecfgdao.Load()))
}

func getResourceCfgCache() *resourceCfgSnapshot {
	v := resourceCfgCache.Load()
	if v == nil {
		return emptyResourceCfg
	}
	cfg, ok := v.(*resourceCfgSnapshot)
	if !ok || cfg == nil {
		return emptyResourceCfg
	}
	return cfg
}

func toResourceCfgSnapshot(row *entity.UploadResourceCfg) *resourceCfgSnapshot {
	if row == nil {
		return emptyResourceCfg
	}
	s := &resourceCfgSnapshot{
		ResourceDomain:   normalizeDomain(row.ResourceDomain),
		DefaultAvatarUrl: strings.TrimSpace(row.DefaultAvatarUrl),
	}
	if s.DefaultAvatarUrl == "" {
		s.DefaultAvatarUrl = defaultAvatarURL
	}
	return s
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return defaultResourceDomain
	}
	lower := strings.ToLower(domain)
	if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
		domain = "http://" + domain
	}
	return strings.TrimRight(domain, "/")
}
