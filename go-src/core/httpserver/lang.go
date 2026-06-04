package httpserver

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/constants/lang"
	"xr-game-server/errercode"
)

// AcceptLanguageHeader 客户端语言协商头(标准 HTTP 头)
const AcceptLanguageHeader = "Accept-Language"

// GetLang 从请求头中解析客户端语言,未传或不识别时回落到默认语言
func GetLang(r *ghttp.Request) lang.Lang {
	return lang.Parse(r.GetHeader(AcceptLanguageHeader))
}

// GetLangFromContext 从上下文关联的请求中解析客户端语言
func GetLangFromContext(ctx context.Context) lang.Lang {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return lang.DefaultLang
	}
	return GetLang(r)
}

// WriteFailJson 将错误码与本地化文字写回客户端(供中间件等无法借助主响应中间件的场景使用)
func WriteFailJson(r *ghttp.Request, code int) {
	resp := CreateFail(code)
	resp.Message = errercode.GetMsg(errercode.XRCode(code), GetLang(r))
	r.Response.WriteJson(resp)
}
