package privacypolicy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/privacypolicycfgdao"
	"xr-game-server/dto/privacypolicydto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetPrivacyPolicyCfg(_ context.Context, _ *privacypolicydto.GetPrivacyPolicyCfgReq) (*privacypolicydto.GetPrivacyPolicyCfgRes, error) {
	cfg := privacypolicycfgdao.Load()
	if cfg == nil {
		return &privacypolicydto.GetPrivacyPolicyCfgRes{Cfg: nil}, nil
	}
	return &privacypolicydto.GetPrivacyPolicyCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SavePrivacyPolicyCfg(_ context.Context, req *privacypolicydto.SavePrivacyPolicyCfgReq) (*privacypolicydto.SavePrivacyPolicyCfgRes, error) {
	url := strings.TrimSpace(req.PrivacyPolicyUrl)
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := privacypolicycfgdao.Load()
	row := &entity.PrivacyPolicyCfg{
		PrivacyPolicyUrl: url,
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
	if err := privacypolicycfgdao.Save(row); err != nil {
		return nil, err
	}
	reloadCfgMemory()
	return &privacypolicydto.SavePrivacyPolicyCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *entity.PrivacyPolicyCfg) *privacypolicydto.PrivacyPolicyCfgItem {
	if cfg == nil {
		return nil
	}
	return &privacypolicydto.PrivacyPolicyCfgItem{
		ID:               strconv.FormatUint(cfg.ID, 10),
		PrivacyPolicyUrl: cfg.PrivacyPolicyUrl,
		CreatedAt:        formatTime(cfg.CreatedAt),
		UpdatedAt:        formatTime(cfg.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
