package preload

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/preloadcfgdto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
	"xr-game-server/module/resourcemonitor"
)

func GetPreloadCfg(_ context.Context, _ *preloadcfgdto.GetPreloadCfgReq) (*preloadcfgdto.GetPreloadCfgRes, error) {
	cfg := cfgdao.LoadPreloadCfg()
	if cfg == nil {
		return &preloadcfgdto.GetPreloadCfgRes{Cfg: defaultPreloadCfgItem()}, nil
	}
	return &preloadcfgdto.GetPreloadCfgRes{Cfg: toPreloadCfgItem(cfg)}, nil
}

func SavePreloadCfg(_ context.Context, req *preloadcfgdto.SavePreloadCfgReq) (*preloadcfgdto.SavePreloadCfgRes, error) {
	existing := cfgdao.LoadPreloadCfg()
	row := &entity.PreloadCfg{
		RecentLoginLimit:                req.RecentLoginLimit,
		HotRestartAuth:                  strings.TrimSpace(req.HotRestartAuth),
		HotRestartFlushTimeout:          req.HotRestartFlushTimeout,
		HotRestartExitTimeout:           req.HotRestartExitTimeout,
		MemoryLimitM:                    req.MemoryLimitM,
		IpGeoDbPath:                     strings.TrimSpace(req.IpGeoDbPath),
		ResourceMetricCollectionEnabled: req.ResourceMetricCollectionEnabled,
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
	resourcemonitor.SetCollectionEnabled(row.ResourceMetricCollectionEnabled)
	return &preloadcfgdto.SavePreloadCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func defaultPreloadCfgItem() *preloadcfgdto.PreloadCfgItem {
	return &preloadcfgdto.PreloadCfgItem{
		RecentLoginLimit:                cfgdao.DefaultRecentLoginPreloadLimit,
		HotRestartAuth:                  cfgdao.DefaultHotRestartAuth,
		HotRestartFlushTimeout:          cfgdao.DefaultHotRestartFlushTimeout,
		HotRestartExitTimeout:           cfgdao.DefaultHotRestartExitTimeout,
		MemoryLimitM:                    cfgdao.DefaultMemoryLimitM,
		IpGeoDbPath:                     cfgdao.DefaultIpGeoDbPath,
		ResourceMetricCollectionEnabled: cfgdao.DefaultResourceMetricCollectionEnabled,
	}
}

func toPreloadCfgItem(cfg *entity.PreloadCfg) *preloadcfgdto.PreloadCfgItem {
	if cfg == nil {
		return nil
	}
	limit := cfg.RecentLoginLimit
	if limit <= 0 {
		limit = cfgdao.DefaultRecentLoginPreloadLimit
	}
	memoryM := cfg.MemoryLimitM
	if memoryM <= 0 {
		memoryM = cfgdao.DefaultMemoryLimitM
	}
	auth := strings.TrimSpace(cfg.HotRestartAuth)
	if auth == "" {
		auth = cfgdao.DefaultHotRestartAuth
	}
	flushSec := cfg.HotRestartFlushTimeout
	if flushSec <= 0 {
		flushSec = cfgdao.DefaultHotRestartFlushTimeout
	}
	exitSec := cfg.HotRestartExitTimeout
	if exitSec <= 0 {
		exitSec = cfgdao.DefaultHotRestartExitTimeout
	}
	ipGeoPath := strings.TrimSpace(cfg.IpGeoDbPath)
	if ipGeoPath == "" {
		ipGeoPath = cfgdao.DefaultIpGeoDbPath
	}
	return &preloadcfgdto.PreloadCfgItem{
		ID:                              strconv.FormatUint(cfg.ID, 10),
		RecentLoginLimit:                limit,
		HotRestartAuth:                  auth,
		HotRestartFlushTimeout:          flushSec,
		HotRestartExitTimeout:           exitSec,
		MemoryLimitM:                    memoryM,
		IpGeoDbPath:                     ipGeoPath,
		CreatedAt:                       formatPreloadCfgTime(cfg.CreatedAt),
		ResourceMetricCollectionEnabled: cfg.ResourceMetricCollectionEnabled,
		UpdatedAt:                       formatPreloadCfgTime(cfg.UpdatedAt),
	}
}

func formatPreloadCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
