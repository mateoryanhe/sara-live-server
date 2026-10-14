package recharge

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dao/userloginlocationdao"
	"xr-game-server/entity/recharge"
	"xr-game-server/module/ipgeo"
)

func resolveUserCountryFromContext(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	ip := strings.TrimSpace(r.GetClientIp())
	if ip == "" {
		return ""
	}
	cfCountry := r.Header.Get("CF-IPCountry")
	return ipgeo.ResolveCountryCode(ip, cfCountry)
}

// applyOrderUserCountry 下单时写入用户国家(优先请求 IP 解析,与登录国家逻辑一致; CMS 等无 IP 时回退用户登录/注册国家).
func applyOrderUserCountry(ctx context.Context, order *entity.RechargeOrder) {
	if order == nil {
		return
	}
	code := resolveUserCountryFromContext(ctx)
	if code == "" && order.UserId > 0 {
		if loc := userloginlocationdao.GetByUserId(order.UserId); loc != nil {
			code = strings.TrimSpace(loc.LoginCountry)
			if code == "" {
				code = strings.TrimSpace(loc.RegisterCountry)
			}
		}
	}
	if code != "" {
		order.SetUserCountry(code)
	}
}
