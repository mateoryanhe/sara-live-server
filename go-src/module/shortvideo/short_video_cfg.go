package shortvideo

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"
	"xr-game-server/dao/shortvideocfgdao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const (
	defaultMaxFileSize      uint64 = 100 * 1024 * 1024
	defaultMaxDuration      uint32 = 60
	defaultFreeWatchSeconds uint32 = 7
	defaultEntryEnabled     uint8  = entity.ShortVideoCfgEntryEnabled
)

var shortVideoCfgCache atomic.Value // *shortvideodto.AppShortVideoCfgRes

func reloadShortVideoCfgMemory() {
	cfg := shortvideocfgdao.Get()
	shortVideoCfgCache.Store(toAppShortVideoCfgRes(cfg))
}

func getAppShortVideoCfgCache() *shortvideodto.AppShortVideoCfgRes {
	v := shortVideoCfgCache.Load()
	if v == nil {
		return defaultAppShortVideoCfg()
	}
	cfg, ok := v.(*shortvideodto.AppShortVideoCfgRes)
	if !ok || cfg == nil {
		return defaultAppShortVideoCfg()
	}
	return cfg
}

func defaultAppShortVideoCfg() *shortvideodto.AppShortVideoCfgRes {
	return &shortvideodto.AppShortVideoCfgRes{
		MaxFileSize:      defaultMaxFileSize,
		MaxDuration:      defaultMaxDuration,
		FreeWatchSeconds: defaultFreeWatchSeconds,
		EntryEnabled:     defaultEntryEnabled,
	}
}

func GetShortVideoCfg(_ context.Context, _ *shortvideodto.GetShortVideoCfgReq) (*shortvideodto.GetShortVideoCfgRes, error) {
	cfg := shortvideocfgdao.Get()
	if cfg == nil {
		return &shortvideodto.GetShortVideoCfgRes{Cfg: nil}, nil
	}
	return &shortvideodto.GetShortVideoCfgRes{Cfg: toShortVideoCfgItem(cfg)}, nil
}

func SaveShortVideoCfg(_ context.Context, req *shortvideodto.SaveShortVideoCfgReq) (*shortvideodto.SaveShortVideoCfgRes, error) {
	existing := shortvideocfgdao.Get()
	row := &entity.ShortVideoCfg{
		MaxFileSize:      req.MaxFileSize,
		MaxDuration:      req.MaxDuration,
		FreeWatchSeconds: req.FreeWatchSeconds,
		EntryEnabled:     req.EntryEnabled,
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
	if err := shortvideocfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadShortVideoCfgMemory()
	return &shortvideodto.SaveShortVideoCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func GetAppShortVideoCfg(_ context.Context, _ *shortvideodto.AppShortVideoCfgReq) (*shortvideodto.AppShortVideoCfgRes, error) {
	cfg := getAppShortVideoCfgCache()
	return &shortvideodto.AppShortVideoCfgRes{
		MaxFileSize:      cfg.MaxFileSize,
		MaxDuration:      cfg.MaxDuration,
		FreeWatchSeconds: cfg.FreeWatchSeconds,
		EntryEnabled:     cfg.EntryEnabled,
	}, nil
}

func toShortVideoCfgItem(cfg *entity.ShortVideoCfg) *shortvideodto.ShortVideoCfgItem {
	if cfg == nil {
		return nil
	}
	return &shortvideodto.ShortVideoCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		MaxFileSize:      cfg.MaxFileSize,
		MaxDuration:      cfg.MaxDuration,
		FreeWatchSeconds: cfg.FreeWatchSeconds,
		EntryEnabled:     cfg.EntryEnabled,
		CreatedAt:        formatShortVideoCfgTime(cfg.CreatedAt),
		UpdatedAt:        formatShortVideoCfgTime(cfg.UpdatedAt),
	}
}

func toAppShortVideoCfgRes(cfg *entity.ShortVideoCfg) *shortvideodto.AppShortVideoCfgRes {
	if cfg == nil {
		return defaultAppShortVideoCfg()
	}
	return &shortvideodto.AppShortVideoCfgRes{
		MaxFileSize:      cfg.MaxFileSize,
		MaxDuration:      cfg.MaxDuration,
		FreeWatchSeconds: cfg.FreeWatchSeconds,
		EntryEnabled:     cfg.EntryEnabled,
	}
}

func formatShortVideoCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
