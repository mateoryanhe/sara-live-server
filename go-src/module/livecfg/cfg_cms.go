package livecfg

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/dao/livecfgdao"
	"xr-game-server/dto/livecfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetLiveCfg(_ context.Context, _ *livecfgdto.GetLiveCfgReq) (*livecfgdto.GetLiveCfgRes, error) {
	cfg := livecfgdao.Load()
	if cfg == nil {
		return &livecfgdto.GetLiveCfgRes{Cfg: nil}, nil
	}
	return &livecfgdto.GetLiveCfgRes{Cfg: toLiveCfgItem(cfg)}, nil
}

func SaveLiveCfg(_ context.Context, req *livecfgdto.SaveLiveCfgReq) (*livecfgdto.SaveLiveCfgRes, error) {
	existing := livecfgdao.Load()
	row := &entity.LiveCfg{
		PaidDanmakuPrice: req.PaidDanmakuPrice,
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
	if err := livecfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadLiveCfgMemory()
	return &livecfgdto.SaveLiveCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toLiveCfgItem(cfg *entity.LiveCfg) *livecfgdto.LiveCfgItem {
	if cfg == nil {
		return nil
	}
	return &livecfgdto.LiveCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		PaidDanmakuPrice: cfg.PaidDanmakuPrice,
		CreatedAt:        formatLiveCfgTime(cfg.CreatedAt),
		UpdatedAt:        formatLiveCfgTime(cfg.UpdatedAt),
	}
}

func formatLiveCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
