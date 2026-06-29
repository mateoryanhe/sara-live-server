package httpserver

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	logBodySkipped      = "[文件/图片上传,已省略报文]"
	apiRequestEndLogKey = "apiRequestEndLog"
)

type apiRequestLog struct {
	Phase      string
	Message    string
	ReqId      string
	Url        string
	Ip         string
	AuthId     string
	Body       any
	DurationMs int64
	Code       int
	Response   any
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

func responseBodyForLog(r *ghttp.Request, resp []byte) any {
	if shouldSkipLogBody(r) {
		return logBodySkipped
	}
	return string(resp)
}

func logAPIRequest(ctx context.Context, entry apiRequestLog) {
	if entry.Phase == "start" {
		g.Log().Infof(ctx, "%v,reqId=%v,url=%v,ip=%v,authId=%v,请求数据=%v",
			entry.Message, entry.ReqId, entry.Url, entry.Ip, entry.AuthId, entry.Body)
		return
	}
	if entry.Body != nil {
		g.Log().Infof(ctx, "reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v,请求数据=%v,",
			entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Body)
		return
	}
	if entry.Code != 0 {
		g.Log().Infof(ctx, "reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v，响应错误码数据=%v,",
			entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Code)
		return
	}
	g.Log().Infof(ctx, "reqId=%v,%v,处理时间=%vms,url=%v,ip=%v,authId=%v，响应数据=%v",
		entry.ReqId, entry.Message, entry.DurationMs, entry.Url, entry.Ip, entry.AuthId, entry.Response)
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

func logAPIRequestEnd(r *ghttp.Request, authId, message string, durationMs int64, code int, body, response any) {
	entry := apiRequestLog{
		Phase:      "end",
		Message:    message,
		ReqId:      r.GetHeader(ReqId, ""),
		Url:        r.URL.RequestURI(),
		Ip:         r.GetClientIp(),
		AuthId:     authId,
		DurationMs: durationMs,
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
	AuthId   string
	SysError bool
	Code     int
	Resp     []byte
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
	err := r.GetError()
	if err != nil {
		if pending.SysError {
			logAPIRequestEnd(r, pending.AuthId, "出现无法捕获的错误", durationMs, 0, requestBodyForLog(r), nil)
			return
		}
		logAPIRequestEnd(r, pending.AuthId, "错误码应答", durationMs, pending.Code, nil, nil)
		return
	}

	message := "正常响应"
	if durationMs >= LongDoTime {
		message = "请求处理时间过长"
	}
	logAPIRequestEnd(r, pending.AuthId, message, durationMs, 0, nil, responseBodyForLog(r, pending.Resp))
}

// requestDurationMs 与 GoFrame Access Log 一致: LeaveTime - EnterTime(毫秒)
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
