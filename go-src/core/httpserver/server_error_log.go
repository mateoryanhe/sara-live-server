package httpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"xr-game-server/errercode"
)

func initErrorLogger() {
	ctx := context.Background()
	if g.Cfg().MustGet(ctx, "logger.error").IsEmpty() {
		return
	}
	_ = g.Log("error")
}

func errorLoggerEnabled() bool {
	return !g.Cfg().MustGet(context.Background(), "logger.error").IsEmpty()
}

func shouldSkipHTTPErrorLog(r *ghttp.Request) bool {
	if r == nil {
		return true
	}
	return strings.Contains(r.RequestURI, "/logQuery/")
}

func hookHTTPErrorLogAfterOutput(r *ghttp.Request) {
	if !errorLoggerEnabled() || shouldSkipHTTPErrorLog(r) {
		return
	}
	err := r.GetError()
	if err == nil || errercode.IsBusiness(err) {
		return
	}
	writeHTTPErrorLog(r, err)
}

// writeHTTPErrorLog 写入 HTTP 框架错误日志到 logger.error,格式与 GF Server 内置 error 日志保持一致。
func writeHTTPErrorLog(r *ghttp.Request, err error) {
	if r == nil || err == nil {
		return
	}
	if r.LeaveTime == nil {
		r.LeaveTime = gtime.Now()
	}
	enterTime := r.EnterTime
	if enterTime == nil {
		enterTime = gtime.Now()
	}

	code := gerror.Code(err)
	codeDetailStr := ""
	if codeDetail := code.Detail(); codeDetail != nil {
		codeDetailStr = gstr.Replace(fmt.Sprintf(`%+v`, codeDetail), "\n", " ")
	}
	content := fmt.Sprintf(
		`%d "%s %s %s %s %s" %.3f, %s, "%s", "%s", %d, "%s", "%+v"`,
		r.Response.Status,
		r.Method,
		r.GetSchema(),
		r.Host,
		r.URL.String(),
		r.Proto,
		float64(r.LeaveTime.Sub(enterTime).Milliseconds())/1000,
		r.GetClientIp(),
		r.Referer(),
		r.UserAgent(),
		code.Code(),
		code.Message(),
		codeDetailStr,
	)
	if stack := gerror.Stack(err); stack != "" {
		content += "\nStack:\n" + stack
	} else {
		content += ", " + err.Error()
	}
	g.Log("error").Error(r.Context(), content)
}
