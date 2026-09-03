package recharge

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/yhpaydto"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

func initYhPayCfg() {
	cfgdao.InitYhPayCfgDao()
	cfgdao.ReloadYhPayCfgCache()
}

func GetYhPayCfg(_ context.Context, _ *yhpaydto.GetYhPayCfgReq) (*yhpaydto.GetYhPayCfgRes, error) {
	cfg := cfgdao.GetYhPayCfgCached()
	if cfg == nil {
		return &yhpaydto.GetYhPayCfgRes{Cfg: nil}, nil
	}
	return &yhpaydto.GetYhPayCfgRes{Cfg: toYhPayCfgItem(cfg)}, nil
}

func SaveYhPayCfg(_ context.Context, req *yhpaydto.SaveYhPayCfgReq) (*yhpaydto.SaveYhPayCfgRes, error) {
	merchantCode := strings.TrimSpace(req.MerchantCode)
	apiKey := strings.TrimSpace(req.ApiKey)
	apiHost := strings.TrimRight(strings.TrimSpace(req.ApiHost), "/")
	if merchantCode == "" || apiKey == "" || apiHost == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetYhPayCfgCached()
	row := &entity.YhPayCfg{
		Enabled:         true,
		MerchantCode:    merchantCode,
		ApiKey:          apiKey,
		ApiHost:         apiHost,
		CallbackBaseUrl: strings.TrimRight(strings.TrimSpace(req.CallbackBaseUrl), "/"),
		ReturnUrl:       strings.TrimSpace(req.ReturnUrl),
		FailedReturnUrl: strings.TrimSpace(req.FailedReturnUrl),
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
		row.CryptoApiHost = existing.CryptoApiHost
		row.CryptoNetwork = existing.CryptoNetwork
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
		row.CryptoApiHost = existing.CryptoApiHost
		row.CryptoNetwork = existing.CryptoNetwork
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveYhPayCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadYhPayCfgCache()
	return &yhpaydto.SaveYhPayCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toYhPayCfgItem(cfg *entity.YhPayCfg) *yhpaydto.YhPayCfgItem {
	if cfg == nil {
		return nil
	}
	return &yhpaydto.YhPayCfgItem{
		ID:              strconv.FormatUint(cfg.ID, 10),
		Enabled:         cfg.Enabled,
		MerchantCode:    cfg.MerchantCode,
		ApiKey:          cfg.ApiKey,
		ApiHost:         cfg.ApiHost,
		CallbackBaseUrl: cfg.CallbackBaseUrl,
		ReturnUrl:       cfg.ReturnUrl,
		FailedReturnUrl: cfg.FailedReturnUrl,
		CreatedAt:       formatYhPayTime(cfg.CreatedAt),
		UpdatedAt:       formatYhPayTime(cfg.UpdatedAt),
	}
}

func formatYhPayTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
