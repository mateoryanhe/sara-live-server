package privacypolicy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/privacypolicydto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetPrivacyPolicyCfg(_ context.Context, _ *privacypolicydto.GetPrivacyPolicyCfgReq) (*privacypolicydto.GetPrivacyPolicyCfgRes, error) {
	cfg := cfgdao.LoadPrivacyPolicyCfg()
	if cfg == nil {
		return &privacypolicydto.GetPrivacyPolicyCfgRes{Cfg: nil}, nil
	}
	return &privacypolicydto.GetPrivacyPolicyCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SavePrivacyPolicyCfg(_ context.Context, req *privacypolicydto.SavePrivacyPolicyCfgReq) (*privacypolicydto.SavePrivacyPolicyCfgRes, error) {
	url := strings.TrimSpace(req.PrivacyPolicyUrl)
	termsUrl := strings.TrimSpace(req.TermsOfServiceUrl)
	creatorTermsUrl := strings.TrimSpace(req.CreatorTermsUrl)
	roomOwnerTermsUrl := strings.TrimSpace(req.RoomOwnerTermsUrl)
	vipDescUrl := strings.TrimSpace(req.VipDescUrl)
	aboutSiteUrl := strings.TrimSpace(req.AboutSiteUrl)
	safetyCenterUrl := strings.TrimSpace(req.SafetyCenterUrl)
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if termsUrl != "" && !strings.HasPrefix(termsUrl, "http://") && !strings.HasPrefix(termsUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if creatorTermsUrl != "" && !strings.HasPrefix(creatorTermsUrl, "http://") && !strings.HasPrefix(creatorTermsUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if roomOwnerTermsUrl != "" && !strings.HasPrefix(roomOwnerTermsUrl, "http://") && !strings.HasPrefix(roomOwnerTermsUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if vipDescUrl != "" && !strings.HasPrefix(vipDescUrl, "http://") && !strings.HasPrefix(vipDescUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if aboutSiteUrl != "" && !strings.HasPrefix(aboutSiteUrl, "http://") && !strings.HasPrefix(aboutSiteUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if safetyCenterUrl != "" && !strings.HasPrefix(safetyCenterUrl, "http://") && !strings.HasPrefix(safetyCenterUrl, "https://") {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.LoadPrivacyPolicyCfg()
	row := &entity.PrivacyPolicyCfg{
		PrivacyPolicyUrl:  url,
		TermsOfServiceUrl: termsUrl,
		CreatorTermsUrl:   creatorTermsUrl,
		RoomOwnerTermsUrl: roomOwnerTermsUrl,
		VipDescUrl:        vipDescUrl,
		AboutSiteUrl:      aboutSiteUrl,
		SafetyCenterUrl:   safetyCenterUrl,
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
	if err := cfgdao.SavePrivacyPolicyCfg(row); err != nil {
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
		ID:                strconv.FormatUint(cfg.ID, 10),
		PrivacyPolicyUrl:  cfg.PrivacyPolicyUrl,
		TermsOfServiceUrl: cfg.TermsOfServiceUrl,
		CreatorTermsUrl:   cfg.CreatorTermsUrl,
		RoomOwnerTermsUrl: cfg.RoomOwnerTermsUrl,
		VipDescUrl:        cfg.VipDescUrl,
		AboutSiteUrl:      cfg.AboutSiteUrl,
		SafetyCenterUrl:   cfg.SafetyCenterUrl,
		CreatedAt:         formatTime(cfg.CreatedAt),
		UpdatedAt:         formatTime(cfg.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
