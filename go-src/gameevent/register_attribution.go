package gameevent

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	userentity "xr-game-server/entity/user"
)

// RegisterAttribution 注册时 Meta CAPI 高级匹配参数(来自请求与 deviceInfo).
type RegisterAttribution struct {
	ClientIPAddress string
	ClientUserAgent string
	MetaFbp         string
	MetaFbc         string
}

func registerAttributionFromCtx(ctx context.Context, deviceInfo *userentity.DeviceInfo) RegisterAttribution {
	var out RegisterAttribution
	r := g.RequestFromCtx(ctx)
	if r != nil {
		out.ClientIPAddress = strings.TrimSpace(r.GetClientIp())
		out.ClientUserAgent = strings.TrimSpace(r.Header.Get("User-Agent"))
		fbp, fbc := metaCookiesFromHeader(r.Header.Get("Cookie"))
		out.MetaFbp = fbp
		out.MetaFbc = fbc
	}
	if deviceInfo != nil {
		if v := strings.TrimSpace(deviceInfo.MetaFbp); v != "" {
			out.MetaFbp = v
		}
		if v := strings.TrimSpace(deviceInfo.MetaFbc); v != "" {
			out.MetaFbc = v
		}
	}
	return out
}

func metaCookiesFromHeader(cookieHeader string) (fbp, fbc string) {
	cookieHeader = strings.TrimSpace(cookieHeader)
	if cookieHeader == "" {
		return "", ""
	}
	for _, part := range strings.Split(cookieHeader, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		name, val, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		val = strings.TrimSpace(val)
		switch name {
		case "_fbp":
			fbp = val
		case "_fbc":
			fbc = val
		}
	}
	return fbp, fbc
}
