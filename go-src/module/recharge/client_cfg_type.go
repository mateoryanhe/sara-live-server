package recharge

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/recharge"
)

const clientPlatformHeader = "X-Platform"

// AppRechargeCfgTypeFromClient 按 X-Platform 与 User-Agent 选择充值档位 cfgType。
// X-Platform=web → 渠道(3)；ios → iOS(1)；android → Google(2)；未传时回退 UA 推断原生端。
func AppRechargeCfgTypeFromClient(platform, ua string) uint8 {
	switch normalizeClientPlatform(platform) {
	case "web":
		return entity.RechargeCfgTypeChannel
	case "ios", "apple":
		return entity.RechargeCfgTypeIOS
	case "android":
		return entity.RechargeCfgTypeGoogle
	}
	return AppRechargeCfgTypeFromUserAgent(ua)
}

func normalizeClientPlatform(platform string) string {
	return strings.ToLower(strings.TrimSpace(platform))
}

// AppRechargeCfgTypeFromUserAgent 无 X-Platform 时按 User-Agent 推断原生 iOS / Android。
func AppRechargeCfgTypeFromUserAgent(ua string) uint8 {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return entity.RechargeCfgTypeGoogle
	}
	lower := strings.ToLower(ua)
	if isIOSNativeUserAgent(lower) {
		return entity.RechargeCfgTypeIOS
	}
	if isAndroidNativeUserAgent(lower) {
		return entity.RechargeCfgTypeGoogle
	}
	return entity.RechargeCfgTypeGoogle
}

func isIOSNativeUserAgent(lowerUA string) bool {
	if strings.Contains(lowerUA, "cfnetwork") {
		return true
	}
	if !strings.Contains(lowerUA, "darwin") {
		return false
	}
	return strings.Contains(lowerUA, "iphone") ||
		strings.Contains(lowerUA, "ipad") ||
		strings.Contains(lowerUA, "ipod")
}

func isAndroidNativeUserAgent(lowerUA string) bool {
	if strings.Contains(lowerUA, "okhttp") || strings.Contains(lowerUA, "dalvik") {
		return true
	}
	return strings.Contains(lowerUA, "android") && !strings.Contains(lowerUA, "mozilla/")
}

func clientPlatformFromContext(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.GetHeader(clientPlatformHeader))
}

func userAgentFromContext(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.UserAgent())
}

func appRechargeCfgTypeFromContext(ctx context.Context) uint8 {
	return AppRechargeCfgTypeFromClient(clientPlatformFromContext(ctx), userAgentFromContext(ctx))
}
