package httpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

func initAccessLogger() {
	ctx := context.Background()
	if g.Cfg().MustGet(ctx, "logger.access").IsEmpty() {
		return
	}
	_ = g.Log("access")
}

func accessLoggerEnabled() bool {
	return !g.Cfg().MustGet(context.Background(), "logger.access").IsEmpty()
}

func shouldSkipAccessLog(r *ghttp.Request) bool {
	if r == nil {
		return true
	}
	return strings.Contains(r.RequestURI, "/logQuery/")
}

func hookAccessLogAfterOutput(r *ghttp.Request) {
	if !accessLoggerEnabled() || shouldSkipAccessLog(r) {
		return
	}
	writeAccessLog(r)
}

// writeAccessLog 写入 access 日志,格式与 GF Server 内置 access 日志保持一致。
func writeAccessLog(r *ghttp.Request) {
	if r == nil {
		return
	}
	if r.LeaveTime == nil {
		r.LeaveTime = gtime.Now()
	}
	enterTime := r.EnterTime
	if enterTime == nil {
		enterTime = gtime.Now()
	}
	content := fmt.Sprintf(
		`%d "%s %s %s %s %s" %.3f, %s, "%s", "%s"`,
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
	)
	g.Log("access").Print(r.Context(), content)
}
