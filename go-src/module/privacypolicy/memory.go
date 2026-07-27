package privacypolicy

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity"
)

type cfgSnapshot struct {
	PrivacyPolicyUrl  string
	TermsOfServiceUrl string
}

var (
	cfgCache         atomic.Value // *cfgSnapshot
	emptyCfgSnapshot = &cfgSnapshot{}
)

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadPrivacyPolicyCfg()))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return emptyCfgSnapshot
	}
	cfg, ok := v.(*cfgSnapshot)
	if !ok || cfg == nil {
		return emptyCfgSnapshot
	}
	return cfg
}

func toCfgSnapshot(row *entity.PrivacyPolicyCfg) *cfgSnapshot {
	if row == nil {
		return emptyCfgSnapshot
	}
	return &cfgSnapshot{
		PrivacyPolicyUrl:  row.PrivacyPolicyUrl,
		TermsOfServiceUrl: row.TermsOfServiceUrl,
	}
}

func GetPrivacyPolicyUrl() string {
	return getCfgCache().PrivacyPolicyUrl
}

func GetTermsOfServiceUrl() string {
	return getCfgCache().TermsOfServiceUrl
}
