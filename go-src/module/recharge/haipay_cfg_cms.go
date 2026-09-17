package recharge

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/haipaydto"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

func initHaiPayCfg() {
	cfgdao.InitHaiPayCfgDao()
	cfgdao.ReloadHaiPayCfgCache()
	RegisterChannelPayProvider(&haiPayProvider{})
}

func GetHaiPayCfg(_ context.Context, _ *haipaydto.GetHaiPayCfgReq) (*haipaydto.GetHaiPayCfgRes, error) {
	cfg := cfgdao.GetHaiPayCfgCached()
	if cfg == nil {
		return &haipaydto.GetHaiPayCfgRes{Cfg: nil}, nil
	}
	return &haipaydto.GetHaiPayCfgRes{Cfg: toHaiPayCfgItem(cfg)}, nil
}

func SaveHaiPayCfg(_ context.Context, req *haipaydto.SaveHaiPayCfgReq) (*haipaydto.SaveHaiPayCfgRes, error) {
	apiHost := strings.TrimRight(strings.TrimSpace(req.ApiHost), "/")
	merchantSecret := strings.TrimSpace(req.MerchantSecretKey)
	merchantPrivate := strings.TrimSpace(req.MerchantPrivateKey)
	if apiHost == "" || req.GlobalCashierAppId <= 0 || merchantSecret == "" || merchantPrivate == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.PayoutEnabled && strings.TrimSpace(req.CallbackBaseUrl) == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetHaiPayCfgCached()
	row := &entity.HaiPayCfg{
		ApiHost:            apiHost,
		GlobalCashierAppId: req.GlobalCashierAppId,
		TVisable:           req.TVisable,
		MerchantSecretKey:  merchantSecret,
		MerchantPrivateKey: merchantPrivate,
		CallbackBaseUrl:    strings.TrimRight(strings.TrimSpace(req.CallbackBaseUrl), "/"),
		ReturnUrl:          strings.TrimSpace(req.ReturnUrl),
		FailReturnUrl:      strings.TrimSpace(req.FailReturnUrl),
		CancelUrl:          strings.TrimSpace(req.CancelUrl),
		PaymentMethods:     strings.TrimSpace(req.PaymentMethods),
		Subject:            strings.TrimSpace(req.Subject),
		PayoutEnabled:      req.PayoutEnabled,
		PayoutAppIds:       strings.TrimSpace(req.PayoutAppIds),
		PayoutSubject:      strings.TrimSpace(req.PayoutSubject),
	}
	if row.Subject == "" {
		row.Subject = "Recharge"
	}
	if row.PayoutSubject == "" {
		row.PayoutSubject = "GuildSettlement"
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
	if err := cfgdao.SaveHaiPayCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadHaiPayCfgCache()
	return &haipaydto.SaveHaiPayCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toHaiPayCfgItem(cfg *entity.HaiPayCfg) *haipaydto.HaiPayCfgItem {
	if cfg == nil {
		return nil
	}
	return &haipaydto.HaiPayCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		ApiHost:            cfg.ApiHost,
		GlobalCashierAppId: cfg.GlobalCashierAppId,
		TVisable:           cfg.TVisable,
		MerchantSecretKey:  cfg.MerchantSecretKey,
		MerchantPrivateKey: cfg.MerchantPrivateKey,
		CallbackBaseUrl:    cfg.CallbackBaseUrl,
		ReturnUrl:          cfg.ReturnUrl,
		FailReturnUrl:      cfg.FailReturnUrl,
		CancelUrl:          cfg.CancelUrl,
		PaymentMethods:     cfg.PaymentMethods,
		Subject:            cfg.Subject,
		PayoutEnabled:      cfg.PayoutEnabled,
		PayoutAppIds:       cfg.PayoutAppIds,
		PayoutSubject:      cfg.PayoutSubject,
		CreatedAt:          formatHaiPayCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatHaiPayCfgTime(cfg.UpdatedAt),
	}
}

func formatHaiPayCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
