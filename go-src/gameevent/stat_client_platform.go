package gameevent

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

const clientPlatformHeader = "X-Platform"

// NormalizeStatClientPlatform 分钟活跃统计用平台维度：android / ios / web，其余为空(未知/未传).
func NormalizeStatClientPlatform(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "android":
		return "android"
	case "ios":
		return "ios"
	case "web":
		return "web"
	default:
		return ""
	}
}

// StatClientPlatformFromCtx 从 HTTP 请求头 X-Platform 读取并规范化.
func StatClientPlatformFromCtx(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return NormalizeStatClientPlatform(r.GetHeader(clientPlatformHeader))
}
