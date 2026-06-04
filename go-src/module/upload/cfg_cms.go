package upload

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/uploadresourcecfgdao"
	"xr-game-server/dto/uploaddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetUploadResourceCfg(_ context.Context, _ *uploaddto.GetUploadResourceCfgReq) (*uploaddto.GetUploadResourceCfgRes, error) {
	cfg := uploadresourcecfgdao.Load()
	if cfg == nil {
		return &uploaddto.GetUploadResourceCfgRes{Cfg: nil}, nil
	}
	return &uploaddto.GetUploadResourceCfgRes{Cfg: toUploadResourceCfgItem(cfg)}, nil
}

func SaveUploadResourceCfg(_ context.Context, req *uploaddto.SaveUploadResourceCfgReq) (*uploaddto.SaveUploadResourceCfgRes, error) {
	existing := uploadresourcecfgdao.Load()
	row := &entity.UploadResourceCfg{
		ResourceDomain:   strings.TrimSpace(req.ResourceDomain),
		DefaultAvatarUrl: strings.TrimSpace(req.DefaultAvatarUrl),
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
	if err := uploadresourcecfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadResourceCfgMemory()
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
		ID:               strconv.FormatUint(cfg.ID, 10),
		ResourceDomain:   snap.ResourceDomain,
		DefaultAvatarUrl: snap.DefaultAvatarUrl,
		CreatedAt:        formatCfgTime(cfg.CreatedAt),
		UpdatedAt:        formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
