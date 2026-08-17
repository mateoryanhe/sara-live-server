package auth

import (
	"strings"

	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/accountcfg"
)

// 常见模拟器/虚拟机 CPU 或硬件标识(小写匹配)
var simulatorCPUKeywords = []string{
	"goldfish",   // Android Emulator
	"ranchu",     // Android Emulator (QEMU)
	"vbox86",     // VirtualBox based emulator
	"ttvm",       // some Android emulators
	"genymotion", // Genymotion
	"qemu",       // QEMU
	"virtualbox",
	"bluestacks",
	"nox",
	"memu",
	"ldplayer",
	"mumuvm",
	"android_x86",
}

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

// isSimulatorByCPU 根据 CPU 型号/硬件标识判断是否模拟器
func isSimulatorByCPU(info *entity.DeviceInfo) bool {
	if info == nil {
		return false
	}
	cpu := strings.ToLower(strings.TrimSpace(info.CpuModel))
	if cpu == "" {
		return false
	}
	deviceType := strings.ToLower(strings.TrimSpace(info.DeviceType))
	// iOS 真机为 ARM;上报 x86/i386 基本可判定为模拟器
	if deviceType == "ios" {
		if cpu == "x86_64" || cpu == "i386" || strings.Contains(cpu, "simulator") {
			return true
		}
	}
	for _, kw := range simulatorCPUKeywords {
		if strings.Contains(cpu, kw) {
			return true
		}
	}
	return false
}
