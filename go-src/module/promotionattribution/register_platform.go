package promotionattribution

import (
	"context"
	"strings"
)

const (
	RegisterPlatformH5 = "h5"
	RegisterPlatformApp = "app"
)

// 与 auth.DeviceChannel / H5DeviceChannel / EmailChannel 保持一致(避免 import auth 循环依赖).
const (
	registerPlatformDeviceChannel uint = 4
	registerPlatformH5DeviceChannel uint = 6
	registerPlatformEmailChannel    uint = 7
)

// ResolveRegisterPromotionPlatform 仅按 accounts.channel 记录注册端(不读 HTTP Header).
func ResolveRegisterPromotionPlatform(_ context.Context, accountChannel uint) string {
	switch accountChannel {
	case registerPlatformH5DeviceChannel, registerPlatformEmailChannel:
		return RegisterPlatformH5
	case registerPlatformDeviceChannel:
		return RegisterPlatformApp
	default:
		return ""
	}
}

func IsRegisterPlatformH5(platform string) bool {
	return strings.TrimSpace(platform) == RegisterPlatformH5
}

func IsRegisterPlatformApp(platform string) bool {
	return strings.TrimSpace(platform) == RegisterPlatformApp
}
