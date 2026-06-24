package shortvideo

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
)

// 配置短视频大小,计费等
const (
	defaultMaxFileSize      uint64 = 100 * 1024 * 1024
	defaultMaxCoverFileSize uint32 = 5
	defaultMaxDuration      uint32 = 60
	defaultFreeWatchSeconds uint32 = 8
	defaultEntryEnabled     uint8  = entity.ShortVideoCfgEntryEnabled
)

func defaultAppShortVideoCfg() *shortvideodto.AppShortVideoCfgRes {
	return &shortvideodto.AppShortVideoCfgRes{
		MaxFileSize:      defaultMaxFileSize,
		MaxCoverFileSize: defaultMaxCoverFileSize,
		MaxDuration:      defaultMaxDuration,
		FreeWatchSeconds: defaultFreeWatchSeconds,
		EntryEnabled:     defaultEntryEnabled,
	}
}

func GetShortVideoCfg(_ context.Context, _ *shortvideodto.GetShortVideoCfgReq) (*shortvideodto.GetShortVideoCfgRes, error) {
	cfg := shortvideodao.Get()
	if cfg == nil {
		return &shortvideodto.GetShortVideoCfgRes{Cfg: &shortvideodto.ShortVideoCfgItem{
			MaxFileSize:      defaultMaxFileSize,
			MaxCoverFileSize: defaultMaxCoverFileSize,
			MaxDuration:      defaultMaxDuration,
			FreeWatchSeconds: defaultFreeWatchSeconds,
			EntryEnabled:     defaultEntryEnabled,
		}}, nil
	}
	return &shortvideodto.GetShortVideoCfgRes{Cfg: toShortVideoCfgItem(cfg)}, nil
}

func SaveShortVideoCfg(_ context.Context, req *shortvideodto.SaveShortVideoCfgReq) (*shortvideodto.SaveShortVideoCfgRes, error) {
	existing := shortvideodao.Get()
	row := &entity.ShortVideoCfg{
		MaxFileSize:      req.MaxFileSize,
		MaxCoverFileSize: req.MaxCoverFileSize,
		MaxDuration:      req.MaxDuration,
		FreeWatchSeconds: req.FreeWatchSeconds,
		EntryEnabled:     req.EntryEnabled,
	}
	if existing != nil {
		row.ID = existing.ID
		row.UpdatedAt = existing.UpdatedAt
	} else {
		row.CreatedAt = time.Now()
		row.UpdatedAt = time.Now()
		row.ID = snowflake.GetId()
	}

	if err := shortvideodao.Save(row); err != nil {
		return nil, err
	}
	return &shortvideodto.SaveShortVideoCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func GetAppShortVideoCfg(_ context.Context, _ *shortvideodto.AppShortVideoCfgReq) (*shortvideodto.AppShortVideoCfgRes, error) {
	cfg := shortvideodao.Get()
	if cfg == nil {
		return defaultAppShortVideoCfg(), nil
	}
	return &shortvideodto.AppShortVideoCfgRes{
		MaxFileSize:      cfg.MaxFileSize,
		MaxCoverFileSize: getShortVideoMaxCoverFileSize(),
		MaxDuration:      cfg.MaxDuration,
		FreeWatchSeconds: cfg.FreeWatchSeconds,
		EntryEnabled:     cfg.EntryEnabled,
	}, nil
}

func toShortVideoCfgItem(cfg *entity.ShortVideoCfg) *shortvideodto.ShortVideoCfgItem {
	if cfg == nil {
		return nil
	}
	maxCoverFileSize := cfg.MaxCoverFileSize
	if maxCoverFileSize == 0 {
		maxCoverFileSize = defaultMaxCoverFileSize
	}
	return &shortvideodto.ShortVideoCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		MaxFileSize:      cfg.MaxFileSize,
		MaxCoverFileSize: maxCoverFileSize,
		MaxDuration:      cfg.MaxDuration,
		FreeWatchSeconds: cfg.FreeWatchSeconds,
		EntryEnabled:     cfg.EntryEnabled,
		CreatedAt:        formatShortVideoCfgTime(cfg.CreatedAt),
		UpdatedAt:        formatShortVideoCfgTime(cfg.UpdatedAt),
	}
}

func formatShortVideoCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func getShortVideoMaxFileSize() uint64 {
	cfg := shortvideodao.Get()
	if cfg == nil || cfg.MaxFileSize == 0 {
		return defaultMaxFileSize
	}
	return cfg.MaxFileSize
}

func getShortVideoMaxCoverFileSize() uint32 {
	cfg := shortvideodao.Get()
	if cfg == nil || cfg.MaxCoverFileSize == 0 {
		return defaultMaxCoverFileSize
	}
	return cfg.MaxCoverFileSize
}

func getShortVideoMaxDuration() uint32 {
	cfg := shortvideodao.Get()
	if cfg == nil || cfg.MaxDuration == 0 {
		return defaultMaxDuration
	}
	return cfg.MaxDuration
}

func getShortVideoFreeWatchSeconds() uint32 {
	cfg := shortvideodao.Get()
	if cfg == nil {
		return defaultFreeWatchSeconds
	}
	return cfg.FreeWatchSeconds
}
