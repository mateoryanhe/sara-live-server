package httpserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	logBodySkipped       = "[文件/图片上传,已省略报文]"
	logRespSkippedBySize = "[响应体已省略,大小=%v字节]"
	apiRequestEndLogKey  = "apiRequestEndLog"
)

type apiRequestLog struct {
	Phase         string
	Message       string
	ReqId         string
	Url           string
	Ip            string
	AuthId        string
	Body          any
	DurationMs    int64
	PreHandlerMs  int64
	HandlerMs     int64
	SerializeMs   int64
	RespBytes     int
	Code          int
	Response      any
	Err           error
	Stack         string
	SlowBreakdown bool
}

// shouldSkipLogBody 请求体为 multipart 文件/图片上传时,日志不输出原始报文
func shouldSkipLogBody(r *ghttp.Request) bool {
	if r == nil {
		return false
	}
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.HasPrefix(ct, "multipart/form-data") {
		return true
	}
	if strings.HasPrefix(ct, "image/") {
		return true
	}
	if strings.HasPrefix(ct, "application/octet-stream") {
		return true
	}
	if r.GetUploadFile("file") != nil {
		return true
	}
	return false
}

func requestBodyForLog(r *ghttp.Request) any {
	if shouldSkipLogBody(r) {
		return logBodySkipped
	}
	return r.GetBodyString()
}

func responseBodyForLog(r *ghttp.Request, resp []byte, slowBreakdown bool) any {
	if shouldSkipLogBody(r) {
		return logBodySkipped
	}
	if slowBreakdown {
		return fmt.Sprintf(logRespSkippedBySize, len(resp))
	}
	return string(resp)
}

func slowTimingSuffix(entry apiRequestLog) string {
	if !entry.SlowBreakdown {
		return ""
	}
	return fmt.Sprintf(",preHandlerMs=%vms,handlerMs=%vms,serializeMs=%vms,respBytes=%v",
		entry.PreHandlerMs, entry.HandlerMs, entry.SerializeMs, entry.RespBytes)
}

func logTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func logAPIRequest(ctx context.Context, entry apiRequestLog) {
	now := logTimeNow()
	if entry.Phase == "start" {
		g.Log().Infof(ctx, "time=%v,%v,reqId=%v,url=%v,ip=%v,authId=%v,请求数据=%v",
			now, entry.Message, entry.ReqId, entry.Url, entry.Ip, entry.AuthId, entry.Body)
		return
	}
	if entry.Body != nil {
		if entry.Err != nil {
			g.Log().Infof(ctx, "time=%v,reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v,请求数据=%v,错误信息=%v",
				now, entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Body, entry.Err)
			if entry.Stack != "" {
				g.Log().Errorf(ctx, "time=%v,reqId=%v,url=%v,堆栈信息=%v", now, entry.ReqId, entry.Url, entry.Stack)
			}
			return
		}
		g.Log().Infof(ctx, "time=%v,reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v,请求数据=%v,",
			now, entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Body)
		return
	}
	if entry.Code != 0 {
		if entry.Err != nil {
			g.Log().Infof(ctx, "time=%v,reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v,响应错误码=%v,错误信息=%v,",
				now, entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Code, entry.Err)
			return
		}
		g.Log().Infof(ctx, "time=%v,reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v,响应错误码=%v,",
			now, entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Code)
		return
	}
	g.Log().Infof(ctx, "time=%v,reqId=%v,%v,处理时间=%vms%s,url=%v,ip=%v,authId=%v,响应数据=%v",
		now, entry.ReqId, entry.Message, entry.DurationMs, slowTimingSuffix(entry), entry.Url, entry.Ip, entry.AuthId, entry.Response)
}

func logAPIRequestStart(r *ghttp.Request, authId, message string) {
	logAPIRequest(r.Context(), apiRequestLog{
		Phase:   "start",
		Message: message,
		ReqId:   r.GetHeader(ReqId, ""),
		Url:     r.URL.RequestURI(),
		Ip:      r.GetClientIp(),
		AuthId:  authId,
		Body:    requestBodyForLog(r),
	})
}

func logAPIRequestEnd(r *ghttp.Request, authId, message string, durationMs, preHandlerMs, handlerMs, serializeMs int64, respBytes int, slowBreakdown bool, code int, body, response any, err error, stack string) {
	entry := apiRequestLog{
		Phase:         "end",
		Message:       message,
		ReqId:         r.GetHeader(ReqId, ""),
		Url:           r.URL.RequestURI(),
		Ip:            r.GetClientIp(),
		AuthId:        authId,
		DurationMs:    durationMs,
		PreHandlerMs:  preHandlerMs,
		HandlerMs:     handlerMs,
		SerializeMs:   serializeMs,
		RespBytes:     respBytes,
		SlowBreakdown: slowBreakdown,
		Err:           err,
		Stack:         stack,
	}
	if code != 0 {
		entry.Code = code
	}
	if body != nil {
		entry.Body = body
	}
	if response != nil {
		entry.Response = response
	}
	logAPIRequest(r.Context(), entry)
}

type apiRequestEndPending struct {
	AuthId       string
	SysError     bool
	Code         int
	Resp         []byte
	PreHandlerMs int64
	HandlerMs    int64
	SerializeMs  int64
	RespBytes    int
}

func stashAPIRequestEndLog(r *ghttp.Request, pending *apiRequestEndPending) {
	r.SetCtxVar(apiRequestEndLogKey, pending)
}

func hookAPIRequestEndLog(r *ghttp.Request) {
	v := r.GetCtxVar(apiRequestEndLogKey)
	if v.IsNil() {
		return
	}
	pending, ok := v.Val().(*apiRequestEndPending)
	if !ok || pending == nil {
		return
	}
	if r.LeaveTime == nil {
		r.LeaveTime = gtime.Now()
	}
	durationMs := requestDurationMs(r)
	slowBreakdown := durationMs >= LongDoTime
	err := r.GetError()
	if err != nil {
		if pending.SysError {
			logAPIRequestEnd(r, pending.AuthId, "出现无法捕获的错误", durationMs,
				pending.PreHandlerMs, pending.HandlerMs, pending.SerializeMs, pending.RespBytes, slowBreakdown,
				0, requestBodyForLog(r), nil, err, gerror.Stack(err))
			return
		}
		logAPIRequestEnd(r, pending.AuthId, "错误码应答", durationMs,
			pending.PreHandlerMs, pending.HandlerMs, pending.SerializeMs, pending.RespBytes, slowBreakdown,
			pending.Code, nil, nil, err, "")
		return
	}

	message := "正常响应"
	if slowBreakdown {
		message = "请求处理时间过长"
	}
	logAPIRequestEnd(r, pending.AuthId, message, durationMs,
		pending.PreHandlerMs, pending.HandlerMs, pending.SerializeMs, pending.RespBytes, slowBreakdown,
		0, nil, responseBodyForLog(r, pending.Resp, slowBreakdown), nil, "")
}

// requestDurationMs 从请求进入到当前时刻/LeaveTime 的耗时(毫秒).
// HookAfterOutput 中的总耗时包含:读请求体、鉴权、业务处理、JSON序列化、写响应、gzip 等.
func requestDurationMs(r *ghttp.Request) int64 {
	if r == nil || r.EnterTime == nil {
		return -1
	}
	if r.LeaveTime != nil {

		return r.LeaveTime.Sub(r.EnterTime).Milliseconds()
	}

	return gtime.Now().Sub(r.EnterTime).Milliseconds()
}

func authIdFromToken(r *ghttp.Request) string {
	token := r.GetHeader("Authorization", "")
	if token == "" {
		return ""
	}
	return strings.Split(token, ".")[0]
}

// authIdFromRequest App 从 Authorization 解析,CMS 从 authId header 读取
func authIdFromRequest(r *ghttp.Request) string {
	if id := authIdFromToken(r); id != "" {
		return id
	}
	return r.GetHeader(AuthId)
}
