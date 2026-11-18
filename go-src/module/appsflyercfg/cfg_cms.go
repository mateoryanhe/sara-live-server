package appsflyercfg

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/appsflyercfgdto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetAppsFlyerCfg(_ context.Context, _ *appsflyercfgdto.GetAppsFlyerCfgReq) (*appsflyercfgdto.GetAppsFlyerCfgRes, error) {
	cfg := cfgdao.GetAppsFlyerCfgCached()
	if cfg == nil {
		return &appsflyercfgdto.GetAppsFlyerCfgRes{Cfg: nil}, nil
	}
	return &appsflyercfgdto.GetAppsFlyerCfgRes{Cfg: toAppsFlyerCfgItem(cfg)}, nil
}

func SaveAppsFlyerCfg(_ context.Context, req *appsflyercfgdto.SaveAppsFlyerCfgReq) (*appsflyercfgdto.SaveAppsFlyerCfgRes, error) {
	androidDevKey := strings.TrimSpace(req.AndroidDevKey)
	iosDevKey := strings.TrimSpace(req.IosDevKey)
	iosAppId := strings.TrimSpace(req.IosAppId)
	if androidDevKey == "" && iosDevKey == "" && iosAppId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetAppsFlyerCfgCached()
	row := &sysentity.AppsFlyerCfg{
		AndroidDevKey: androidDevKey,
		IosDevKey:     iosDevKey,
		IosAppId:      iosAppId,
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
	now := time.Now()
	row.UpdatedAt = now
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if err := cfgdao.SaveAppsFlyerCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadAppsFlyerCfgCache()
	return &appsflyercfgdto.SaveAppsFlyerCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toAppsFlyerCfgItem(cfg *sysentity.AppsFlyerCfg) *appsflyercfgdto.AppsFlyerCfgItem {
	if cfg == nil {
		return nil
	}
	return &appsflyercfgdto.AppsFlyerCfgItem{
		ID:            strconv.FormatUint(cfg.ID, 10),
		AndroidDevKey: cfg.AndroidDevKey,
		IosDevKey:     cfg.IosDevKey,
		IosAppId:      cfg.IosAppId,
		CreatedAt:     formatAppsFlyerCfgTime(cfg.CreatedAt),
		UpdatedAt:     formatAppsFlyerCfgTime(cfg.UpdatedAt),
	}
}

func formatAppsFlyerCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}