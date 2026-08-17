package auth

import (
	"strings"

	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/accountcfg"
	"xr-game-server/module/simulatorcpukeyword"
)

// ensureSimulatorLoginAllowed 配置开启拦截时,按上报 cpuModel 拒绝模拟器登录(默认不拦截)
func ensureSimulatorLoginAllowed(info *entity.DeviceInfo) error {
	if info == nil || !isSimulatorByCPU(info) {
		return nil
	}
	if !accountcfg.IsSimulatorLoginBlocked() {
		return nil
	}
	return errercode.CreateCode(errercode.SimulatorLoginDenied)
}

// isSimulatorByCPU 根据 CPU 型号/硬件标识判断是否模拟器(CMS 关键词缓存模糊匹配)
func isSimulatorByCPU(info *entity.DeviceInfo) bool {
	if info == nil {
		return false
	}
	cpu := strings.ToLower(strings.TrimSpace(info.CpuModel))
	if cpu == "" {
		return false
	}
	// 真机多为 ARM; 上报含 x86(如 x86/x86_64) 一律视为模拟器(iOS/Android 通用)
	if strings.Contains(cpu, "x86") {
		return true
	}
	return simulatorcpukeyword.MatchCPU(cpu)
}
