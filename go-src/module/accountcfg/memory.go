package accountcfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/cms"
)

type cfgSnapshot struct {
	CancelAccountByCodeEnabled bool
	SimulatorLoginEnabled      bool
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadAccountCfg()))
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{SimulatorLoginEnabled: true}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{SimulatorLoginEnabled: true}
	}
	return snap
}

func toCfgSnapshot(row *entity.AccountCfg) *cfgSnapshot {
	if row == nil {
		return &cfgSnapshot{SimulatorLoginEnabled: true}
	}
	return &cfgSnapshot{
		CancelAccountByCodeEnabled: row.CancelAccountByCodeEnabled,
		SimulatorLoginEnabled:      row.SimulatorLoginEnabled,
	}
}

// IsCancelAccountByCodeEnabled 注销码销户接口是否开启
func IsCancelAccountByCodeEnabled() bool {
	return getCfgCache().CancelAccountByCodeEnabled
}

// IsSimulatorLoginEnabled 模拟器是否允许登录
func IsSimulatorLoginEnabled() bool {
	return getCfgCache().SimulatorLoginEnabled
}
