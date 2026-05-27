package agora

import (
	"context"
	"strconv"
	"time"
	"xr-game-server/dao/agoracfgdao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetAgoraCfg(_ context.Context, _ *agoradto.GetAgoraCfgReq) (*agoradto.GetAgoraCfgRes, error) {
	cfg := agoracfgdao.Load()
	if cfg == nil {
		return &agoradto.GetAgoraCfgRes{Cfg: nil}, nil
	}
	return &agoradto.GetAgoraCfgRes{Cfg: toAgoraCfgItem(cfg)}, nil
}

func SaveAgoraCfg(_ context.Context, req *agoradto.SaveAgoraCfgReq) (*agoradto.SaveAgoraCfgRes, error) {
	existing := agoracfgdao.Load()
	row := &entity.AgoraCfg{
		AppId:              req.AppId,
		AppCertificate:     req.AppCertificate,
		TokenExpireSeconds: req.TokenExpireSeconds,
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
	if err := agoracfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadAgoraCfgMemory()
	return &agoradto.SaveAgoraCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toAgoraCfgItem(cfg *entity.AgoraCfg) *agoradto.AgoraCfgItem {
	if cfg == nil {
		return nil
	}
	return &agoradto.AgoraCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		AppId:              cfg.AppId,
		AppCertificate:     cfg.AppCertificate,
		TokenExpireSeconds: cfg.TokenExpireSeconds,
		CreatedAt:          formatAgoraCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatAgoraCfgTime(cfg.UpdatedAt),
	}
}

func formatAgoraCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
