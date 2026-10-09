package metapixel

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/metapixeldto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetMetaPixelCfg(_ context.Context, _ *metapixeldto.GetMetaPixelCfgReq) (*metapixeldto.GetMetaPixelCfgRes, error) {
	cfg := cfgdao.ResolveMetaPixelCfg()
	if cfg == nil {
		return &metapixeldto.GetMetaPixelCfgRes{Cfg: nil}, nil
	}
	return &metapixeldto.GetMetaPixelCfgRes{Cfg: toMetaPixelCfgItem(cfg)}, nil
}

func SaveMetaPixelCfg(_ context.Context, req *metapixeldto.SaveMetaPixelCfgReq) (*metapixeldto.SaveMetaPixelCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	pixelId := strings.TrimSpace(req.PixelId)
	accessToken := strings.TrimSpace(req.AccessToken)
	testEventCode := strings.TrimSpace(req.TestEventCode)
	if req.Enabled == 1 {
		if pixelId == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if accessToken == "" {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
	}
	cfgdao.EnsureMetaPixelCfgDao()
	existing := cfgdao.ResolveMetaPixelCfg()
	row := &sysentity.MetaPixelCfg{
		Enabled:       req.Enabled,
		PixelId:       pixelId,
		AccessToken:   accessToken,
		TestEventCode: testEventCode,
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
	if err := cfgdao.SaveMetaPixelCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadMetaPixelCfgCache()
	OnMetaPixelCfgSaved(row)
	return &metapixeldto.SaveMetaPixelCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toMetaPixelCfgItem(cfg *sysentity.MetaPixelCfg) *metapixeldto.MetaPixelCfgItem {
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	return &metapixeldto.MetaPixelCfgItem{
		ID:            strconv.FormatUint(cfg.ID, 10),
		Enabled:       cfg.Enabled,
		PixelId:       cfg.PixelId,
		AccessToken:   cfg.AccessToken,
		TestEventCode: cfg.TestEventCode,
		CreatedAt:     formatMetaPixelTime(cfg.CreatedAt),
		UpdatedAt:     formatMetaPixelTime(cfg.UpdatedAt),
	}
}

func formatMetaPixelTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func getCfgCached() *sysentity.MetaPixelCfg {
	return cfgdao.ResolveMetaPixelCfg()
}
