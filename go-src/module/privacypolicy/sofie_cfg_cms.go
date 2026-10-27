package privacypolicy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/privacypolicydto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetSofiePrivacyPolicyCfg(_ context.Context, _ *privacypolicydto.GetSofiePrivacyPolicyCfgReq) (*privacypolicydto.GetSofiePrivacyPolicyCfgRes, error) {
	cfg := cfgdao.LoadSofiePrivacyPolicyCfg()
	if cfg == nil {
		return &privacypolicydto.GetSofiePrivacyPolicyCfgRes{Cfg: nil}, nil
	}
	return &privacypolicydto.GetSofiePrivacyPolicyCfgRes{Cfg: toSofieCfgItem(cfg)}, nil
}

func SaveSofiePrivacyPolicyCfg(_ context.Context, req *privacypolicydto.SaveSofiePrivacyPolicyCfgReq) (*privacypolicydto.SaveSofiePrivacyPolicyCfgRes, error) {
	apiBase, err := normalizeApiBase(req.ApiBase)
	if err != nil {
		return nil, err
	}

	url := strings.TrimSpace(req.PrivacyPolicyUrl)
	termsUrl := strings.TrimSpace(req.TermsOfServiceUrl)
	creatorTermsUrl := strings.TrimSpace(req.CreatorTermsUrl)
	roomOwnerTermsUrl := strings.TrimSpace(req.RoomOwnerTermsUrl)
	vipDescUrl := strings.TrimSpace(req.VipDescUrl)
	aboutSiteUrl := strings.TrimSpace(req.AboutSiteUrl)
	safetyCenterUrl := strings.TrimSpace(req.SafetyCenterUrl)

	resolve := func(raw string) (string, error) {
		return resolvePolicyUrl(apiBase, raw)
	}
	if url, err = resolve(url); err != nil {
		return nil, err
	}
	if termsUrl, err = resolve(termsUrl); err != nil {
		return nil, err
	}
	if creatorTermsUrl, err = resolve(creatorTermsUrl); err != nil {
		return nil, err
	}
	if roomOwnerTermsUrl, err = resolve(roomOwnerTermsUrl); err != nil {
		return nil, err
	}
	if vipDescUrl, err = resolve(vipDescUrl); err != nil {
		return nil, err
	}
	if aboutSiteUrl, err = resolve(aboutSiteUrl); err != nil {
		return nil, err
	}
	if safetyCenterUrl, err = resolve(safetyCenterUrl); err != nil {
		return nil, err
	}

	packageName := strings.TrimSpace(req.PackageName)

	existing := cfgdao.LoadSofiePrivacyPolicyCfg()
	row := &entity.SofiePrivacyPolicyCfg{
		PackageName:       packageName,
		ApiBase:           apiBase,
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
	if err := cfgdao.SaveSofiePrivacyPolicyCfg(row); err != nil {
		return nil, err
	}
	reloadSofieCfgMemory()
	return &privacypolicydto.SaveSofiePrivacyPolicyCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toSofieCfgItem(cfg *entity.SofiePrivacyPolicyCfg) *privacypolicydto.PrivacyPolicyCfgItem {
	if cfg == nil {
		return nil
	}
	apiBase := strings.TrimRight(strings.TrimSpace(cfg.ApiBase), "/")
	return &privacypolicydto.PrivacyPolicyCfgItem{
		ID:                strconv.FormatUint(cfg.ID, 10),
		PackageName:       strings.TrimSpace(cfg.PackageName),
		ApiBase:           apiBase,
		PrivacyPolicyUrl:  stripApiBase(apiBase, cfg.PrivacyPolicyUrl),
		TermsOfServiceUrl: stripApiBase(apiBase, cfg.TermsOfServiceUrl),
		CreatorTermsUrl:   stripApiBase(apiBase, cfg.CreatorTermsUrl),
		RoomOwnerTermsUrl: stripApiBase(apiBase, cfg.RoomOwnerTermsUrl),
		VipDescUrl:        stripApiBase(apiBase, cfg.VipDescUrl),
		AboutSiteUrl:      stripApiBase(apiBase, cfg.AboutSiteUrl),
		SafetyCenterUrl:   stripApiBase(apiBase, cfg.SafetyCenterUrl),
		CreatedAt:         formatTime(cfg.CreatedAt),
		UpdatedAt:         formatTime(cfg.UpdatedAt),
	}
}
