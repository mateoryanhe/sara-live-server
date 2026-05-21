package httpserver

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

const logBodySkipped = "[文件/图片上传,已省略报文]"

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
