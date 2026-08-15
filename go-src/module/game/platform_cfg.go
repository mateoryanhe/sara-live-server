package game

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity/game"
	"xr-game-server/errercode"
)

func initGamePlatformCfg() {
	cfgdao.InitGamePlatformCfgDao()
	cfgdao.ReloadGamePlatformCfgCache()
}

func GetGamePlatformCfg(_ context.Context, _ *gameplatformdto.GetGamePlatformCfgReq) (*gameplatformdto.GetGamePlatformCfgRes, error) {
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil {
		return &gameplatformdto.GetGamePlatformCfgRes{Cfg: nil}, nil
	}
	return &gameplatformdto.GetGamePlatformCfgRes{Cfg: toGamePlatformCfgItem(cfg)}, nil
}

func SaveGamePlatformCfg(_ context.Context, req *gameplatformdto.SaveGamePlatformCfgReq) (*gameplatformdto.SaveGamePlatformCfgRes, error) {
	vendorUrl := normalizeVendorUrl(req.VendorUrl)
	token := strings.TrimSpace(req.Token)
	secretKey := strings.TrimSpace(req.SecretKey)
	if vendorUrl == "" || token == "" || secretKey == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetGamePlatformCfgFromMemory()
	row := &entity.GamePlatformCfg{
		VendorUrl: vendorUrl,
		Token:     token,
		SecretKey: secretKey,
		IconUrl:   strings.TrimRight(strings.TrimSpace(req.IconUrl), "/"),
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
	if err := cfgdao.SaveGamePlatformCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGamePlatformCfgCache()
	return &gameplatformdto.SaveGamePlatformCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toGamePlatformCfgItem(cfg *entity.GamePlatformCfg) *gameplatformdto.GamePlatformCfgItem {
	if cfg == nil {
		return nil
	}
	return &gameplatformdto.GamePlatformCfgItem{
		ID:        strconv.FormatUint(cfg.ID, 10),
		VendorUrl: GetVendorUrl(cfg),
		Token:     cfg.Token,
		SecretKey: cfg.SecretKey,
		IconUrl:   cfg.IconUrl,
		CreatedAt: formatGamePlatformCfgTime(cfg.CreatedAt),
		UpdatedAt: formatGamePlatformCfgTime(cfg.UpdatedAt),
	}
}

// GetVendorUrl 获取厂家 API 根地址,未配置时返回默认值.
func GetVendorUrl(cfg *entity.GamePlatformCfg) string {
	if cfg == nil {
		return entity.GamePlatformDefaultVendorUrl
	}
	vendorUrl := normalizeVendorUrl(cfg.VendorUrl)
	if vendorUrl == "" {
		return entity.GamePlatformDefaultVendorUrl
	}
	return vendorUrl
}

func normalizeVendorUrl(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}

func formatGamePlatformCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
