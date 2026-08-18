package liverevenuesharecfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/liverevenuesharecfgdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

const (
	DefaultAnchorSharePercent = 30
	DefaultGuildSharePercent  = 10
)

func GetLiveRevenueShareCfg(_ context.Context, _ *liverevenuesharecfgdto.GetLiveRevenueShareCfgReq) (*liverevenuesharecfgdto.GetLiveRevenueShareCfgRes, error) {
	cfg := cfgdao.GetLiveRevenueShareCfgCached()
	if cfg == nil {
		return &liverevenuesharecfgdto.GetLiveRevenueShareCfgRes{
			Cfg: &liverevenuesharecfgdto.LiveRevenueShareCfgItem{
				AnchorSharePercent: DefaultAnchorSharePercent,
				GuildSharePercent:  DefaultGuildSharePercent,
			},
		}, nil
	}
	return &liverevenuesharecfgdto.GetLiveRevenueShareCfgRes{Cfg: toLiveRevenueShareCfgItem(cfg)}, nil
}

func SaveLiveRevenueShareCfg(_ context.Context, req *liverevenuesharecfgdto.SaveLiveRevenueShareCfgReq) (*liverevenuesharecfgdto.SaveLiveRevenueShareCfgRes, error) {
	if req.AnchorSharePercent < 0 || req.AnchorSharePercent > 100 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.GuildSharePercent < 0 || req.GuildSharePercent > 100 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetLiveRevenueShareCfgCached()
	row := &entity.LiveRevenueShareCfg{
		AnchorSharePercent: req.AnchorSharePercent,
		GuildSharePercent:  req.GuildSharePercent,
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
	if err := cfgdao.SaveLiveRevenueShareCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadLiveRevenueShareCfgCache()
	return &liverevenuesharecfgdto.SaveLiveRevenueShareCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toLiveRevenueShareCfgItem(cfg *entity.LiveRevenueShareCfg) *liverevenuesharecfgdto.LiveRevenueShareCfgItem {
	if cfg == nil {
		return nil
	}
	return &liverevenuesharecfgdto.LiveRevenueShareCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		AnchorSharePercent: cfg.AnchorSharePercent,
		GuildSharePercent:  cfg.GuildSharePercent,
		CreatedAt:          formatCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
