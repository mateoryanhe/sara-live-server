package recharge

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/googleplaydto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func initGooglePlayCfg() {
	cfgdao.InitGooglePlayCfgDao()
	cfgdao.ReloadGooglePlayCfgCache()
}

func GetGooglePlayCfg(_ context.Context, _ *googleplaydto.GetGooglePlayCfgReq) (*googleplaydto.GetGooglePlayCfgRes, error) {
	cfg := cfgdao.GetGooglePlayCfgFromMemory()
	if cfg == nil {
		return &googleplaydto.GetGooglePlayCfgRes{Cfg: nil}, nil
	}
	return &googleplaydto.GetGooglePlayCfgRes{Cfg: toGooglePlayCfgItem(cfg)}, nil
}

func SaveGooglePlayCfg(_ context.Context, req *googleplaydto.SaveGooglePlayCfgReq) (*googleplaydto.SaveGooglePlayCfgRes, error) {
	packageName := strings.TrimSpace(req.PackageName)
	serviceAccountJson := strings.TrimSpace(req.ServiceAccountJson)
	rtdnAudience := strings.TrimSpace(req.RtdnAudience)
	if req.Enabled {
		if packageName == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if serviceAccountJson == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if !json.Valid([]byte(serviceAccountJson)) {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
	}
	if err := validateOptionalURL(rtdnAudience); err != nil {
		return nil, err
	}

	existing := cfgdao.GetGooglePlayCfgFromMemory()
	row := &entity.GooglePlayCfg{
		Enabled:            req.Enabled,
		PackageName:        packageName,
		ServiceAccountJson: serviceAccountJson,
		RtdnAudience:       rtdnAudience,
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
	if err := cfgdao.SaveGooglePlayCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGooglePlayCfgCache()
	ResetGooglePlayClient()
	return &googleplaydto.SaveGooglePlayCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func validateOptionalURL(url string) error {
	if url == "" {
		return nil
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return nil
	}
	return errercode.CreateCode(errercode.InvalidParam)
}

func toGooglePlayCfgItem(cfg *entity.GooglePlayCfg) *googleplaydto.GooglePlayCfgItem {
	if cfg == nil {
		return nil
	}
	return &googleplaydto.GooglePlayCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		Enabled:            cfg.Enabled,
		PackageName:        cfg.PackageName,
		ServiceAccountJson: cfg.ServiceAccountJson,
		RtdnAudience:       cfg.RtdnAudience,
		CreatedAt:          formatGooglePlayCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatGooglePlayCfgTime(cfg.UpdatedAt),
	}
}

func formatGooglePlayCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
