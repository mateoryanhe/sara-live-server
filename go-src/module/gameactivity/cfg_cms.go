package gameactivity

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameactivitydto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetGameActivityCfg(_ context.Context, _ *gameactivitydto.GetGameActivityCfgReq) (*gameactivitydto.GetGameActivityCfgRes, error) {
	cfg := cfgdao.ResolveGameActivityCfg()
	return &gameactivitydto.GetGameActivityCfgRes{Cfg: toGameActivityCfgItem(cfg)}, nil
}

func SaveGameActivityCfg(_ context.Context, req *gameactivitydto.SaveGameActivityCfgReq) (*gameactivitydto.SaveGameActivityCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &sysentity.GameActivityCfg{
		ButtonIcon:  strings.TrimSpace(req.ButtonIcon),
		GameListUrl: strings.TrimSpace(req.GameListUrl),
	}
	cfgdao.EnsureGameActivityCfgDao()
	existing := cfgdao.GetGameActivityCfgFromMemory()
	if existing == nil {
		existing = cfgdao.ResolveGameActivityCfg()
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
	if err := cfgdao.SaveGameActivityCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameActivityCfgCache()
	return &gameactivitydto.SaveGameActivityCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toGameActivityCfgItem(cfg *sysentity.GameActivityCfg) *gameactivitydto.GameActivityCfgItem {
	if cfg == nil || cfg.ID == 0 {
		return nil
	}
	return &gameactivitydto.GameActivityCfgItem{
		ID:          strconv.FormatUint(cfg.ID, 10),
		ButtonIcon:  cfg.ButtonIcon,
		GameListUrl: cfg.GameListUrl,
		CreatedAt:   formatGameActivityTime(cfg.CreatedAt),
		UpdatedAt:   formatGameActivityTime(cfg.UpdatedAt),
	}
}

func formatGameActivityTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
