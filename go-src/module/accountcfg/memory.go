package accountcfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/cms"
)

type cfgSnapshot struct {
	CancelAccountByCodeEnabled bool
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadAccountCfg()))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{}
	}
	return snap
}

func toCfgSnapshot(row *entity.AccountCfg) *cfgSnapshot {
	if row == nil {
		return &cfgSnapshot{}
	}
	return &cfgSnapshot{
		CancelAccountByCodeEnabled: row.CancelAccountByCodeEnabled,
	}
}

// IsCancelAccountByCodeEnabled 注销码销户接口是否开启
func IsCancelAccountByCodeEnabled() bool {
	return getCfgCache().CancelAccountByCodeEnabled
}
