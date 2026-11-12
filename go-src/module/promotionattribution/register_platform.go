package promotionattribution

import (
	"strings"
)

const (
	RegisterPlatformH5  = "h5"
	RegisterPlatformApp = "app"
)

// 与 auth.DeviceChannel / H5DeviceChannel / EmailChannel 保持一致(避免 import auth 循环依赖).
const (
	registerPlatformDeviceChannel   uint = 4
	registerPlatformH5DeviceChannel uint = 6
	registerPlatformEmailChannel    uint = 7
)

// MapClientXPlatformToRegisterPlatform 按请求头 X-Platform 映射客户端报表端:
// web → h5; android/ios → app; 其它空.
func MapClientXPlatformToRegisterPlatform(xPlatform string) string {
	switch strings.ToLower(strings.TrimSpace(xPlatform)) {
	case "web":
		return RegisterPlatformH5
	case "android", "ios":
		return RegisterPlatformApp
	default:
		return ""
	}
}

// ResolveRegisterPromotionPlatformFromChannel 无 X-Platform 时的渠道兜底(历史兼容).
func ResolveRegisterPromotionPlatformFromChannel(accountChannel uint) string {
	switch accountChannel {
	case registerPlatformH5DeviceChannel, registerPlatformEmailChannel:
		return RegisterPlatformH5
	case registerPlatformDeviceChannel:
		return RegisterPlatformApp
	default:
		return ""
	}
}

// ResolveRegisterPromotionPlatform 优先 X-Platform,空则回落 accounts.channel.
func ResolveRegisterPromotionPlatform(xPlatform string, accountChannel uint) string {
	if p := MapClientXPlatformToRegisterPlatform(xPlatform); p != "" {
		return p
	}
	return ResolveRegisterPromotionPlatformFromChannel(accountChannel)
}

func IsRegisterPlatformH5(platform string) bool {
	return strings.TrimSpace(platform) == RegisterPlatformH5
}

func IsRegisterPlatformApp(platform string) bool {
	return strings.TrimSpace(platform) == RegisterPlatformApp
}
