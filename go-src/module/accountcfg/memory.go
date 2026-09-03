package accountcfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/cms"
)

type cfgSnapshot struct {
	CancelAccountByCodeEnabled bool
	BlockSimulatorLogin        bool
	EnvType                    uint8
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
		BlockSimulatorLogin:        row.BlockSimulatorLogin,
		EnvType:                    row.EnvType,
	}
}

// IsCancelAccountByCodeEnabled 注销码销户接口是否开启
func IsCancelAccountByCodeEnabled() bool {
	return getCfgCache().CancelAccountByCodeEnabled
}

// IsSimulatorLoginBlocked 是否拦截模拟器登录(默认 false=不拦截)
func IsSimulatorLoginBlocked() bool {
	return getCfgCache().BlockSimulatorLogin
}

// GetEnvType 环境类型(0正式服,1提审服,2测试服)
func GetEnvType() uint8 {
	return getCfgCache().EnvType
}
