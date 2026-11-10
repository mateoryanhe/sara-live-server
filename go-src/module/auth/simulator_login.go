package auth

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"xr-game-server/core/event"
	"xr-game-server/dao/accountdao"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/accountcfg"
	"xr-game-server/module/simulatorcpukeyword"
	"xr-game-server/module/simulatordevicewhitelist"
)

const (
	SimulatorLoginTypeDevice       = "device"
	SimulatorLoginTypeEmail        = "email"
	SimulatorLoginTypePhone        = "phone"
	SimulatorLoginTypePhoneReg     = "phoneRegister"
	SimulatorLoginTypeFirebase     = "firebase"
	SimulatorLoginTypeCoinMerchant = "coinMerchant"
)

// SimulatorDenyLogOpt 模拟器拦截落库附加信息(设备码/邮箱登录必填)
type SimulatorDenyLogOpt struct {
	LoginType  string
	OpenIdHint string
	Channel    uint
	Email      string
}

// ensureSimulatorLoginAllowed 配置开启拦截时,按上报 cpuModel 拒绝模拟器登录(默认不拦截;设备白名单放行)
func ensureSimulatorLoginAllowed(ctx context.Context, info *entity.DeviceInfo, opt *SimulatorDenyLogOpt) error {
	if info == nil || !isSimulatorByCPU(info) {
		return nil
	}
	if !accountcfg.IsSimulatorLoginBlocked() {
		return nil
	}
	if simulatordevicewhitelist.IsWhitelisted(info.DeviceId) {
		return nil
	}
	publishSimulatorLoginDenied(ctx, info, opt)
	return errercode.CreateCode(errercode.SimulatorLoginDenied)
}

func publishSimulatorLoginDenied(ctx context.Context, info *entity.DeviceInfo, opt *SimulatorDenyLogOpt) {
	if opt == nil {
		return
	}
	// 仅记录设备码快捷登录与 App 邮箱登录
	if opt.LoginType != SimulatorLoginTypeDevice && opt.LoginType != SimulatorLoginTypeEmail {
		return
	}
	ip := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = strings.TrimSpace(r.GetClientIp())
	}
	email := strings.TrimSpace(opt.Email)
	accountId := uint64(0)
	var registerAt *time.Time
	openId := strings.TrimSpace(opt.OpenIdHint)
	if openId != "" && opt.Channel > 0 {
		if account := accountdao.FindActiveAccount(openId, opt.Channel); account != nil {
			accountId = account.ID
			if !account.CreatedAt.IsZero() {
				t := account.CreatedAt
				registerAt = &t
			}
		}
	}
	if accountId == 0 && opt.LoginType == SimulatorLoginTypeEmail && email != "" {
		if account := findActiveAccountByBoundEmail(email); account != nil {
			accountId = account.ID
			if !account.CreatedAt.IsZero() {
				t := account.CreatedAt
				registerAt = &t
			}
		}
	}
	event.Pub(gameevent.SimulatorLoginDeniedEvent, gameevent.NewSimulatorLoginDeniedEventData(
		ip, info, opt.LoginType, email, accountId, registerAt,
	))
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
