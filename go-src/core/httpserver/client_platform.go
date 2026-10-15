package httpserver

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

const ClientPlatformHeader = "X-Platform"

// ClientPlatformFromContext 读取 App/Web 客户端平台标识(小写 trim)。
func ClientPlatformFromContext(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(r.GetHeader(ClientPlatformHeader)))
}

// IsWebClientPlatform X-Platform=web 时为 true。
func IsWebClientPlatform(ctx context.Context) bool {
	return ClientPlatformFromContext(ctx) == "web"
}
