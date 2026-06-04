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
	ResourceDomain                 string
	DefaultAvatarUrl               string
	ImageModerationEnabled         bool
	ImageModerationAccessKeyId     string
	ImageModerationAccessKeySecret string
	ImageModerationRegionId        string
	ImageModerationEndpoint        string
	ImageModerationService         string
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
		ResourceDomain:                 normalizeDomain(row.ResourceDomain),
		DefaultAvatarUrl:               strings.TrimSpace(row.DefaultAvatarUrl),
		ImageModerationEnabled:         row.ImageModerationEnabled,
		ImageModerationAccessKeyId:     strings.TrimSpace(row.ImageModerationAccessKeyId),
		ImageModerationAccessKeySecret: strings.TrimSpace(row.ImageModerationAccessKeySecret),
		ImageModerationRegionId:        strings.TrimSpace(row.ImageModerationRegionId),
		ImageModerationEndpoint:        strings.TrimSpace(row.ImageModerationEndpoint),
		ImageModerationService:         strings.TrimSpace(row.ImageModerationService),
	}
	if s.DefaultAvatarUrl == "" {
		s.DefaultAvatarUrl = defaultAvatarURL
	}
	if s.ImageModerationRegionId == "" {
		s.ImageModerationRegionId = defaultImageModerationRegion
	}
	if s.ImageModerationEndpoint == "" {
		s.ImageModerationEndpoint = defaultImageModerationEndpoint
	}
	if s.ImageModerationService == "" {
		s.ImageModerationService = defaultImageModerationService
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
