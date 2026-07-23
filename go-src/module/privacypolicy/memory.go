package privacypolicy

import (
	"sync/atomic"

	"xr-game-server/dao/privacypolicycfgdao"
	"xr-game-server/entity"
)

type cfgSnapshot struct {
	PrivacyPolicyUrl string
}

var (
	cfgCache         atomic.Value // *cfgSnapshot
	emptyCfgSnapshot = &cfgSnapshot{}
)

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(privacypolicycfgdao.Load()))
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
		PrivacyPolicyUrl: row.PrivacyPolicyUrl,
	}
}

func GetPrivacyPolicyUrl() string {
	return getCfgCache().PrivacyPolicyUrl
}
