package accountcfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/cms"
)

const (
	defaultDeviceAccountMaxCount  = 3
	defaultDeviceCancelDailyLimit = 1
)

type cfgSnapshot struct {
	CancelAccountByCodeEnabled bool
	BlockSimulatorLogin        bool
	EnvType                    uint8
	DeviceRegisterRiskEnabled  bool
	DeviceAccountMaxCount      int
	DeviceCancelDailyLimit     int
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadAccountCfg()))
}

func defaultRiskSnapshot() *cfgSnapshot {
	return &cfgSnapshot{
		DeviceRegisterRiskEnabled:  true,
		DeviceAccountMaxCount:      defaultDeviceAccountMaxCount,
		DeviceCancelDailyLimit:     defaultDeviceCancelDailyLimit,
	}
}

func getCfgCache() *cfgSnapshot {
	v := cfgCache.Load()
	if v == nil {
		return defaultRiskSnapshot()
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return defaultRiskSnapshot()
	}
	return snap
}

func toCfgSnapshot(row *entity.AccountCfg) *cfgSnapshot {
	if row == nil {
		return defaultRiskSnapshot()
	}
	maxCount := row.DeviceAccountMaxCount
	if maxCount <= 0 {
		maxCount = defaultDeviceAccountMaxCount
	}
	dailyLimit := row.DeviceCancelDailyLimit
	if dailyLimit <= 0 {
		dailyLimit = defaultDeviceCancelDailyLimit
	}
	return &cfgSnapshot{
		CancelAccountByCodeEnabled: row.CancelAccountByCodeEnabled,
		BlockSimulatorLogin:        row.BlockSimulatorLogin,
		EnvType:                    row.EnvType,
		DeviceRegisterRiskEnabled:  row.DeviceRegisterRiskEnabled,
		DeviceAccountMaxCount:      maxCount,
		DeviceCancelDailyLimit:     dailyLimit,
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

// IsDeviceRegisterRiskEnabled 设备码注册风控是否开启
func IsDeviceRegisterRiskEnabled() bool {
	return getCfgCache().DeviceRegisterRiskEnabled
}

// GetDeviceAccountMaxCount 同设备最大账号数(含已注销)
func GetDeviceAccountMaxCount() int {
	return getCfgCache().DeviceAccountMaxCount
}

// GetDeviceCancelDailyLimit 同设备每天最多注销次数
func GetDeviceCancelDailyLimit() int {
	return getCfgCache().DeviceCancelDailyLimit
}
