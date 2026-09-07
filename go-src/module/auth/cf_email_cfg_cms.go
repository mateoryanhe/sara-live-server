package auth

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/cfemaildto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func initCfEmailCfg() {
	cfgdao.InitCfEmailCfgDao()
	cfgdao.ReloadCfEmailCfgCache()
}

func GetCfEmailCfg(_ context.Context, _ *cfemaildto.GetCfEmailCfgReq) (*cfemaildto.GetCfEmailCfgRes, error) {
	cfg := cfgdao.GetCfEmailCfgCached()
	if cfg == nil {
		return &cfemaildto.GetCfEmailCfgRes{Cfg: nil}, nil
	}
	return &cfemaildto.GetCfEmailCfgRes{Cfg: toCfEmailCfgItem(cfg)}, nil
}

func SaveCfEmailCfg(_ context.Context, req *cfemaildto.SaveCfEmailCfgReq) (*cfemaildto.SaveCfEmailCfgRes, error) {
	accountId := strings.TrimSpace(req.AccountId)
	apiToken := strings.TrimSpace(req.ApiToken)
	fromEmail := strings.TrimSpace(req.FromEmail)
	if accountId == "" || apiToken == "" || fromEmail == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetCfEmailCfgCached()
	row := &sysentity.CfEmailCfg{
		Enabled:   req.Enabled,
		AccountId: accountId,
		ApiToken:  apiToken,
		FromEmail: fromEmail,
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
	if err := cfgdao.SaveCfEmailCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadCfEmailCfgCache()
	return &cfemaildto.SaveCfEmailCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfEmailCfgItem(cfg *sysentity.CfEmailCfg) *cfemaildto.CfEmailCfgItem {
	if cfg == nil {
		return nil
	}
	return &cfemaildto.CfEmailCfgItem{
		ID:        strconv.FormatUint(cfg.ID, 10),
		Enabled:   cfg.Enabled,
		AccountId: cfg.AccountId,
		ApiToken:  cfg.ApiToken,
		FromEmail: cfg.FromEmail,
		CreatedAt: formatCfEmailTime(cfg.CreatedAt),
		UpdatedAt: formatCfEmailTime(cfg.UpdatedAt),
	}
}

func formatCfEmailTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
