package effectivelivecfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/effectivelivecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetEffectiveLiveCfg(_ context.Context, _ *effectivelivecfgdto.GetEffectiveLiveCfgReq) (*effectivelivecfgdto.GetEffectiveLiveCfgRes, error) {
	cfg := cfgdao.GetEffectiveLiveCfgCached()
	if cfg == nil {
		return &effectivelivecfgdto.GetEffectiveLiveCfgRes{
			Cfg: &effectivelivecfgdto.EffectiveLiveCfgItem{
				MinSessionMinutes: cfgdao.DefaultEffectiveLiveMinSessionMinutes,
			},
		}, nil
	}
	return &effectivelivecfgdto.GetEffectiveLiveCfgRes{Cfg: toEffectiveLiveCfgItem(cfg)}, nil
}

func SaveEffectiveLiveCfg(_ context.Context, req *effectivelivecfgdto.SaveEffectiveLiveCfgReq) (*effectivelivecfgdto.SaveEffectiveLiveCfgRes, error) {
	if req.MinSessionMinutes <= 0 || req.MinSessionMinutes > cfgdao.MaxEffectiveLiveMinSessionMinutes {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetEffectiveLiveCfgCached()
	row := &entity.EffectiveLiveCfg{MinSessionMinutes: req.MinSessionMinutes}
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
	if err := cfgdao.SaveEffectiveLiveCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadEffectiveLiveCfgCache()
	return &effectivelivecfgdto.SaveEffectiveLiveCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toEffectiveLiveCfgItem(cfg *entity.EffectiveLiveCfg) *effectivelivecfgdto.EffectiveLiveCfgItem {
	if cfg == nil {
		return nil
	}
	return &effectivelivecfgdto.EffectiveLiveCfgItem{
		ID:                strconv.FormatUint(cfg.ID, 10),
		MinSessionMinutes: cfgdao.EffectiveLiveMinSessionMinutes(),
		CreatedAt:         formatCfgTime(cfg.CreatedAt),
		UpdatedAt:         formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
