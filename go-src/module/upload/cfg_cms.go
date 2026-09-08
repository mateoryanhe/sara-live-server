package upload

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/uploaddto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

const uploadCfgMB = 1024 * 1024

func GetUploadResourceCfg(_ context.Context, _ *uploaddto.GetUploadResourceCfgReq) (*uploaddto.GetUploadResourceCfgRes, error) {
	cfg := cfgdao.LoadUploadResourceCfg()
	if cfg == nil {
		return &uploaddto.GetUploadResourceCfgRes{Cfg: nil}, nil
	}
	return &uploaddto.GetUploadResourceCfgRes{Cfg: toUploadResourceCfgItem(cfg)}, nil
}

func SaveUploadResourceCfg(_ context.Context, req *uploaddto.SaveUploadResourceCfgReq) (*uploaddto.SaveUploadResourceCfgRes, error) {
	if req.AppImageMaxSizeMB < 1 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	storagePath := normalizeStoragePath(req.StoragePath)
	if !isAbsoluteStoragePath(storagePath) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadUploadResourceCfg()
	row := &entity.UploadResourceCfg{
		ResourceDomain:                 strings.TrimSpace(req.ResourceDomain),
		StoragePath:                    storagePath,
		CmsExportTtlMinutes:            req.CmsExportTtlMinutes,
		AppImageMaxSize:                uint64(req.AppImageMaxSizeMB) * uploadCfgMB,
		ImageModerationEnabled:         req.ImageModerationEnabled,
		ImageModerationAccessKeyId:     strings.TrimSpace(req.ImageModerationAccessKeyId),
		ImageModerationAccessKeySecret: strings.TrimSpace(req.ImageModerationAccessKeySecret),
		ImageModerationRegionId:        strings.TrimSpace(req.ImageModerationRegionId),
		ImageModerationEndpoint:        strings.TrimSpace(req.ImageModerationEndpoint),
		ImageModerationService:         strings.TrimSpace(req.ImageModerationService),
		S3Enabled:                      req.S3Enabled,
		S3PublicDomain:                 strings.TrimSpace(req.S3PublicDomain),
		S3Endpoint:                     strings.TrimSpace(req.S3Endpoint),
		S3Region:                       defaultS3Region,
		S3Bucket:                       strings.TrimSpace(req.S3Bucket),
		S3AccessKeyId:                  strings.TrimSpace(req.S3AccessKeyId),
		S3SecretAccessKey:              strings.TrimSpace(req.S3SecretAccessKey),
		S3KeyPrefix:                    strings.TrimSpace(req.S3KeyPrefix),
	}
	if row.ImageModerationEnabled {
		if row.ImageModerationAccessKeyId == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing == nil && row.ImageModerationAccessKeySecret == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing != nil && row.ImageModerationAccessKeySecret == "" {
			row.ImageModerationAccessKeySecret = existing.ImageModerationAccessKeySecret
		}
	}
	if row.S3Enabled {
		if row.S3PublicDomain == "" || row.S3Endpoint == "" || row.S3Bucket == "" || row.S3AccessKeyId == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing == nil && row.S3SecretAccessKey == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if existing != nil && row.S3SecretAccessKey == "" {
			row.S3SecretAccessKey = existing.S3SecretAccessKey
		}
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveUploadResourceCfg(row); err != nil {
		return nil, err
	}
	invalidateImageGreenClient()
	invalidateS3Client()
	reloadResourceCfgMemory()
	registerStaticMappings()
	return &uploaddto.SaveUploadResourceCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toUploadResourceCfgItem(cfg *entity.UploadResourceCfg) *uploaddto.UploadResourceCfgItem {
	if cfg == nil {
		return nil
	}
	snap := toResourceCfgSnapshot(cfg)
	return &uploaddto.UploadResourceCfgItem{
		ID:                             strconv.FormatUint(cfg.ID, 10),
		ResourceDomain:                 snap.ResourceDomain,
		StoragePath:                    snap.StoragePath,
		CmsExportTtlMinutes:            snap.CmsExportTtlMinutes,
		AppImageMaxSizeMB:              appImageMaxSizeMB(snap.AppImageMaxSize),
		ImageModerationEnabled:         snap.ImageModerationEnabled,
		ImageModerationAccessKeyId:     cfg.ImageModerationAccessKeyId,
		ImageModerationAccessKeySecret: maskCfgSecret(cfg.ImageModerationAccessKeySecret),
		ImageModerationRegionId:        snap.ImageModerationRegionId,
		ImageModerationEndpoint:        snap.ImageModerationEndpoint,
		ImageModerationService:         snap.ImageModerationService,
		S3Enabled:                      snap.S3Enabled,
		S3PublicDomain:                 snap.S3PublicDomain,
		S3Endpoint:                     snap.S3Endpoint,
		S3Bucket:                       snap.S3Bucket,
		S3AccessKeyId:                  snap.S3AccessKeyId,
		S3SecretAccessKey:              maskCfgSecret(cfg.S3SecretAccessKey),
		S3KeyPrefix:                    strings.TrimSuffix(snap.S3KeyPrefix, "/"),
		CreatedAt:                      formatCfgTime(cfg.CreatedAt),
		UpdatedAt:                      formatCfgTime(cfg.UpdatedAt),
	}
}

func maskCfgSecret(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func appImageMaxSizeMB(size uint64) uint32 {
	if size == 0 {
		size = uploadCfgMB
	}
	mb := size / uploadCfgMB
	if mb < 1 {
		return 1
	}
	return uint32(mb)
}
