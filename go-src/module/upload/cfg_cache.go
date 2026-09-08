package upload

import (
	"path/filepath"
	"strings"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/sys"
)

const (
	defaultResourceDomain   = "http://127.0.0.1"
	defaultStoragePath      = "/home/ec2-user/cdn/images"
	defaultCmsExportTTL     = 30
	defaultAppImageMaxSize  = 1048576 // 1MB
)

type resourceCfgSnapshot struct {
	ResourceDomain                 string
	ResourceDomainConfigured       bool
	StoragePath                    string
	CmsExportTtlMinutes            int
	AppImageMaxSize                uint64
	ImageModerationEnabled         bool
	ImageModerationAccessKeyId     string
	ImageModerationAccessKeySecret string
	ImageModerationRegionId        string
	ImageModerationEndpoint        string
	ImageModerationService         string
	S3Enabled                      bool
	S3PublicDomain                 string
	S3Endpoint                     string
	S3Region                       string
	S3Bucket                       string
	S3AccessKeyId                  string
	S3SecretAccessKey              string
	S3KeyPrefix                    string
}

var (
	resourceCfgCache atomic.Value // *resourceCfgSnapshot
	emptyResourceCfg = &resourceCfgSnapshot{
		ResourceDomain:      defaultResourceDomain,
		StoragePath:         defaultStoragePath,
		CmsExportTtlMinutes: defaultCmsExportTTL,
		AppImageMaxSize:     defaultAppImageMaxSize,
	}
)

func reloadResourceCfgMemory() {
	resourceCfgCache.Store(toResourceCfgSnapshot(cfgdao.LoadUploadResourceCfg()))
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
		ResourceDomainConfigured:       strings.TrimSpace(row.ResourceDomain) != "",
		StoragePath:                    normalizeStoragePath(row.StoragePath),
		CmsExportTtlMinutes:            normalizeCmsExportTtl(row.CmsExportTtlMinutes),
		AppImageMaxSize:                normalizeAppImageMaxSize(row.AppImageMaxSize),
		ImageModerationEnabled:         row.ImageModerationEnabled,
		ImageModerationAccessKeyId:     strings.TrimSpace(row.ImageModerationAccessKeyId),
		ImageModerationAccessKeySecret: strings.TrimSpace(row.ImageModerationAccessKeySecret),
		ImageModerationRegionId:        strings.TrimSpace(row.ImageModerationRegionId),
		ImageModerationEndpoint:        strings.TrimSpace(row.ImageModerationEndpoint),
		ImageModerationService:         strings.TrimSpace(row.ImageModerationService),
		S3Enabled:                      row.S3Enabled,
		S3PublicDomain:                 normalizeOptionalDomain(row.S3PublicDomain),
		S3Endpoint:                     strings.TrimSpace(row.S3Endpoint),
		S3Region:                       strings.TrimSpace(row.S3Region),
		S3Bucket:                       strings.TrimSpace(row.S3Bucket),
		S3AccessKeyId:                  strings.TrimSpace(row.S3AccessKeyId),
		S3SecretAccessKey:              strings.TrimSpace(row.S3SecretAccessKey),
		S3KeyPrefix:                    normalizeS3KeyPrefix(row.S3KeyPrefix),
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
	if s.S3Region == "" {
		s.S3Region = defaultS3Region
	}
	return s
}

func normalizeS3KeyPrefix(prefix string) string {
	prefix = strings.Trim(strings.ReplaceAll(strings.TrimSpace(prefix), "\\", "/"), "/")
	if prefix == "" {
		return ""
	}
	return prefix + "/"
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return defaultResourceDomain
	}
	return normalizeOptionalDomain(domain)
}

// normalizeOptionalDomain 规范化可选域名;空串保持为空(不回落默认本机)
func normalizeOptionalDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return ""
	}
	lower := strings.ToLower(domain)
	if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
		domain = "https://" + domain
	}
	return strings.TrimRight(domain, "/")
}

func normalizeAppImageMaxSize(size uint64) uint64 {
	if size == 0 {
		return defaultAppImageMaxSize
	}
	return size
}

func normalizeStoragePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultStoragePath
	}
	return filepath.Clean(path)
}

func isAbsoluteStoragePath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	if strings.HasPrefix(path, "/") {
		return true
	}
	if len(path) >= 3 && path[1] == ':' {
		c := path[0]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			sep := path[2]
			return sep == '/' || sep == '\\'
		}
	}
	return strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, "//")
}

func normalizeCmsExportTtl(minutes int) int {
	if minutes <= 0 {
		return defaultCmsExportTTL
	}
	return minutes
}

// GetStoragePath 统一文件存储物理目录(头像/图片/短视频/CMS导出)
func GetStoragePath() string {
	return getResourceCfgCache().StoragePath
}

// GetCmsExportTtlMinutes CMS 文件导出过期清理时间(分钟)
func GetCmsExportTtlMinutes() int {
	return getResourceCfgCache().CmsExportTtlMinutes
}

// IsS3Enabled 是否开启 S3 兼容云桶(R2 等)
func IsS3Enabled() bool {
	return getResourceCfgCache().S3Enabled
}

// GetS3PublicDomain 云桶公开访问域名(R2 自定义域);未配置为空
func GetS3PublicDomain() string {
	return getResourceCfgCache().S3PublicDomain
}

// GetFileAccessDomain App/CMS 拼文件 URL 用的域名:开云桶用云桶域,否则用资源域名
func GetFileAccessDomain() string {
	if IsS3Enabled() {
		if d := GetS3PublicDomain(); d != "" {
			return d
		}
	}
	return GetResourceDomain()
}

// GetS3KeyPrefix 环境对象前缀,形如 test/ 或 prod/;未配置为空
func GetS3KeyPrefix() string {
	return getResourceCfgCache().S3KeyPrefix
}
