package httpserver

import (
	"strings"
	"xr-game-server/core/xrjson"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// 日志记录链条
// 收到前端请求 → 鉴权完成 → 读取请求Body → Handler执行完成 → 应答写入缓冲区 → 应答输出完成(gzip+写出)
const (
	logBodyFileSkipped            = "[文件已省略,不输出内容]"
	logBodySkippedNonJSON         = "[非JSON报文,已省略内容]"
	apiResponseBufferWrittenAtKey = "apiResponseBufferWrittenAt"
)

func shouldSkipAPILogChain(r *ghttp.Request) bool {
	if r == nil {
		return true
	}
	return strings.Contains(r.RequestURI, "/logQuery/")
}

func logAPIRequestStart(r *ghttp.Request) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"收到前端请求,enterTime=%v,从队列进入到中间件时间间隔Ms=%vms,method=%v,url=%v,ip=%v,headers=%s",
		requestEnterTimeStr(r),
		waitBeforeMiddlewareMs(r),
		r.Method,
		r.RequestURI,
		r.GetClientIp(),
		requestHeadersForLog(r),
	)
}

func requestEnterTimeStr(r *ghttp.Request) string {
	if r == nil || r.EnterTime == nil {
		return ""
	}
	return r.EnterTime.Time.Format("2006-01-02 15:04:05.000")
}

func waitBeforeMiddlewareMs(r *ghttp.Request) int64 {
	if r == nil || r.EnterTime == nil {
		return -1
	}
	return gtime.Now().Sub(r.EnterTime).Milliseconds()
}

func requestHeadersForLog(r *ghttp.Request) string {
	if r == nil || r.Request == nil {
		return "{}"
	}
	headers := make(map[string][]string, len(r.Request.Header))
	for k, v := range r.Request.Header {
		headers[k] = v
	}
	return string(xrjson.MustMarshal(headers))
}

func logAPIRequestAuth(r *ghttp.Request, authMs int64) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"鉴权完成,reqId=%v,authId=%v,authMs=%vms,url=%v",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		authMs,
		r.RequestURI,
	)
}

func elapsedMs(start *gtime.Time) int64 {
	if start == nil {
		return -1
	}
	return gtime.Now().Sub(start).Milliseconds()
}

func logRequestBodyBeforeHandler(r *ghttp.Request) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	bodyStart := gtime.Now()
	bodyContent, bodyLength := requestBodyContentForLog(r)
	logAPIRequestBody(r, elapsedMs(bodyStart), bodyLength, bodyContent)
}

func requestBodyContentForLog(r *ghttp.Request) (string, int) {
	if r == nil {
		return "", 0
	}
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	if isMultipartRequest(ct) {
		return multipartFormContentForLog(r)
	}
	if strings.Contains(ct, "application/x-www-form-urlencoded") {
		return formContentForLog(r)
	}
	if isBinaryUploadContentType(ct) {
		return logBodyFileSkipped, int(r.ContentLength)
	}

	body := r.GetBodyString()
	bodyLength := len(body)
	if bodyLength == 0 {
		return "", 0
	}
	if isJSONRequestBody(r, body) {
		return body, bodyLength
	}
	return logBodySkippedNonJSON, bodyLength
}

func formContentForLog(r *ghttp.Request) (string, int) {
	formMap := r.GetFormMap()
	if len(formMap) == 0 {
		return "", 0
	}
	content := string(xrjson.MustMarshal(formMap))
	return content, len(content)
}

func multipartFormContentForLog(r *ghttp.Request) (string, int) {
	form := r.GetMultipartForm()
	logMap := make(map[string]any)

	for k, v := range r.GetFormMap() {
		logMap[k] = v
	}
	if form != nil {
		for name, values := range form.Value {
			logMap[name] = multipartValuesForLog(values)
		}
		for name := range form.File {
			logMap[name] = uploadFilesMetaForLog(r.GetUploadFiles(name))
		}
	}
	if len(logMap) == 0 {
		return "", 0
	}
	content := string(xrjson.MustMarshal(logMap))
	bodyLength := int(r.ContentLength)
	if bodyLength <= 0 {
		bodyLength = len(content)
	}
	return content, bodyLength
}

func multipartValuesForLog(values []string) any {
	if len(values) == 0 {
		return ""
	}
	if len(values) == 1 {
		return values[0]
	}
	return values
}

func isMultipartRequest(contentType string) bool {
	return strings.Contains(contentType, "multipart/")
}

func uploadFilesMetaForLog(files ghttp.UploadFiles) any {
	if len(files) == 0 {
		return map[string]any{"skipped": logBodyFileSkipped}
	}
	if len(files) == 1 {
		return uploadFileMetaForLog(files[0])
	}
	metas := make([]map[string]any, 0, len(files))
	for _, file := range files {
		metas = append(metas, uploadFileMetaForLog(file))
	}
	return metas
}

func uploadFileMetaForLog(file *ghttp.UploadFile) map[string]any {
	if file == nil {
		return map[string]any{"skipped": logBodyFileSkipped}
	}
	return map[string]any{
		"filename": file.Filename,
		"size":     file.Size,
		"skipped":  logBodyFileSkipped,
	}
}

func isBinaryUploadContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/") || contentType == "application/octet-stream"
}

func logAPIRequestBody(r *ghttp.Request, bodyMs int64, bodyLength int, bodyContent string) {
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"读取请求Body,reqId=%v,authId=%v,bodyMs=%vms,bodyLength=%v,url=%v,bodyContent=%s",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		bodyMs,
		bodyLength,
		r.RequestURI,
		bodyContent,
	)
}

func logAPIRequestHandler(r *ghttp.Request, handlerMs int64) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"Handler执行完成,reqId=%v,authId=%v,handlerMs=%vms,url=%v",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		handlerMs,
		r.RequestURI,
	)
}

func logAPIRequestResponseWrite(r *ghttp.Request, writeMs int64, respBytes int, respContent string) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"应答序列化,写入框架缓冲区,reqId=%v,authId=%v,writeMs=%vms,respBytes=%v,url=%v,respContent=%s",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		writeMs,
		respBytes,
		r.RequestURI,
		respContent,
	)
}

func stashAPIResponseBufferWrittenAt(r *ghttp.Request) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	r.SetCtxVar(apiResponseBufferWrittenAtKey, gtime.Now())
}

func hookAPIRequestAfterOutput(r *ghttp.Request) {
	if shouldSkipAPILogChain(r) {
		return
	}
	if r == nil {
		return
	}
	v := r.GetCtxVar(apiResponseBufferWrittenAtKey)
	if v.IsNil() {
		return
	}
	writtenAt, ok := v.Val().(*gtime.Time)
	if !ok || writtenAt == nil {
		return
	}
	if r.LeaveTime == nil {
		r.LeaveTime = gtime.Now()
	}
	logAPIRequestAfterOutput(r, elapsedMs(writtenAt))
}

func logAPIRequestAfterOutput(r *ghttp.Request, afterOutputMs int64) {
	if r == nil {
		return
	}
	g.Log().Infof(r.Context(),
		"应答写入到系统缓冲区,输出完成,reqId=%v,authId=%v,afterOutputMs=%vms,gzip=%v,totalMs=%vms,url=%v",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		afterOutputMs,
		strings.EqualFold(r.Response.Header().Get("Content-Encoding"), "gzip"),
		requestDurationMs(r),
		r.RequestURI,
	)
}

func requestDurationMs(r *ghttp.Request) int64 {
	if r == nil || r.EnterTime == nil {
		return -1
	}
	if r.LeaveTime != nil {
		return r.LeaveTime.Sub(r.EnterTime).Milliseconds()
	}
	return gtime.Now().Sub(r.EnterTime).Milliseconds()
}

func isJSONRequestBody(r *ghttp.Request, body string) bool {
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(ct, "application/json") || strings.Contains(ct, "+json") {
		return true
	}
	body = strings.TrimSpace(body)
	if len(body) >= 2 && body[0] == '{' && body[len(body)-1] == '}' {
		return true
	}
	if len(body) >= 2 && body[0] == '[' && body[len(body)-1] == ']' {
		return true
	}
	return false
}
