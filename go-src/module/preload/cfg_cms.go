package preload

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/preloadcfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetPreloadCfg(_ context.Context, _ *preloadcfgdto.GetPreloadCfgReq) (*preloadcfgdto.GetPreloadCfgRes, error) {
	cfg := cfgdao.LoadPreloadCfg()
	if cfg == nil {
		return &preloadcfgdto.GetPreloadCfgRes{Cfg: &preloadcfgdto.PreloadCfgItem{
			RecentLoginLimit: cfgdao.DefaultRecentLoginPreloadLimit,
		}}, nil
	}
	return &preloadcfgdto.GetPreloadCfgRes{Cfg: toPreloadCfgItem(cfg)}, nil
}

func SavePreloadCfg(_ context.Context, req *preloadcfgdto.SavePreloadCfgReq) (*preloadcfgdto.SavePreloadCfgRes, error) {
	existing := cfgdao.LoadPreloadCfg()
	row := &entity.PreloadCfg{
		RecentLoginLimit: req.RecentLoginLimit,
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
	if err := cfgdao.SavePreloadCfg(row); err != nil {
		return nil, err
	}
	return &preloadcfgdto.SavePreloadCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toPreloadCfgItem(cfg *entity.PreloadCfg) *preloadcfgdto.PreloadCfgItem {
	if cfg == nil {
		return nil
	}
	limit := cfg.RecentLoginLimit
	if limit <= 0 {
		limit = cfgdao.DefaultRecentLoginPreloadLimit
	}
	return &preloadcfgdto.PreloadCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		RecentLoginLimit: limit,
		CreatedAt:        formatPreloadCfgTime(cfg.CreatedAt),
		UpdatedAt:        formatPreloadCfgTime(cfg.UpdatedAt),
	}
}

func formatPreloadCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
