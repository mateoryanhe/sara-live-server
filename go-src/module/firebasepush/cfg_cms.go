package firebasepush

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/firebasepushdto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetFirebasePushCfg(_ context.Context, _ *firebasepushdto.GetFirebasePushCfgReq) (*firebasepushdto.GetFirebasePushCfgRes, error) {
	cfg := cfgdao.ResolveFirebasePushCfg()
	return &firebasepushdto.GetFirebasePushCfgRes{Cfg: toFirebasePushCfgItem(cfg)}, nil
}

func SaveFirebasePushCfg(_ context.Context, req *firebasepushdto.SaveFirebasePushCfgReq) (*firebasepushdto.SaveFirebasePushCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	enabled := uint8(0)
	if req.Enabled == 1 {
		enabled = 1
	}
	row := &sysentity.FirebasePushCfg{
		Enabled:          enabled,
		OfflineCallTitle: strings.TrimSpace(req.OfflineCallTitle),
		OfflineCallBody:  strings.TrimSpace(req.OfflineCallBody),
		OfflineCallIcon:  strings.TrimSpace(req.OfflineCallIcon),
	}
	row.ApplyFirebasePushDefaults()

	cfgdao.EnsureFirebasePushCfgDao()
	existing := cfgdao.GetFirebasePushCfgFromMemory()
	if existing == nil {
		existing = loadExistingFirebasePushCfg()
	}
	if existing != nil && existing.ID > 0 {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	} else if req.ID > 0 {
		row.ID = req.ID
	}
	now := time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	if err := cfgdao.SaveFirebasePushCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadFirebasePushCfgCache()
	return &firebasepushdto.SaveFirebasePushCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func loadExistingFirebasePushCfg() *sysentity.FirebasePushCfg {
	cfg := cfgdao.ResolveFirebasePushCfg()
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	return cfg
}

func toFirebasePushCfgItem(cfg *sysentity.FirebasePushCfg) *firebasepushdto.FirebasePushCfgItem {
	if cfg == nil {
		cfg = sysentity.NewDefaultFirebasePushCfg()
	} else {
		cfg.ApplyFirebasePushDefaults()
	}
	id := ""
	if cfg.ID > 0 {
		id = strconv.FormatUint(cfg.ID, 10)
	} else {
		id = "0"
	}
	return &firebasepushdto.FirebasePushCfgItem{
		ID:               id,
		Enabled:          cfg.Enabled,
		OfflineCallTitle: cfg.OfflineCallTitle,
		OfflineCallBody:  cfg.OfflineCallBody,
		OfflineCallIcon:  cfg.OfflineCallIcon,
		CreatedAt:        formatFirebasePushTime(cfg.CreatedAt),
		UpdatedAt:        formatFirebasePushTime(cfg.UpdatedAt),
	}
}

func formatFirebasePushTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
