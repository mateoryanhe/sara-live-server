package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
	"time"
	"xr-game-server/errercode"
)

func apiResponseMiddleware(r *ghttp.Request) {
	r.Middleware.Next()

	userId := ""
	token := r.GetHeader("Authorization", "")
	if token != "" {
		userId = strings.Split(token, ".")[0]
	}
	startTime := time.Now()
	if r.EnterTime != nil {
		startTime = r.EnterTime.Time
	}

	var (
		res   = r.GetHandlerResponse()
		err   = r.GetError()
		ctx   = r.Context()
		param any
	)
	code := errercode.Success
	if err != nil {
		var gErr *errercode.XError
		errors.As(err, &gErr)
		retTime := time.Now().Sub(startTime).Milliseconds()
		if gErr != nil {
			code = int(gErr.Code())
			param = gErr.Param
			g.Log().Infof(ctx, "reqId=%v,错误码应答,处理时间=%vms,url=%v,ip=%v,authId=%v，响应错误码数据=%v,", r.GetHeader(ReqId, ""), retTime, r.URL.RequestURI(), r.GetClientIp(), userId, code)
		} else {
			//清除系统自带的错误日志输出
			r.Response.ClearBuffer()
			code = errercode.SysError
			g.Log().Infof(ctx, "reqId=%v,出现无法捕获的错误,处理时间=%vms,url=%v,ip=%v,authId=%v,请求数据=%v,", r.GetHeader(ReqId, ""), retTime, r.URL.RequestURI(), r.GetClientIp(), userId, requestBodyForLog(r))
		}
	}
	r.Response.Header().Set("Content-Type", contentTypeJson)
	var respData any
	//发现有错误码
	if err != nil {
		failResp := CreateFailAndParam(code, param)
		failResp.Message = errercode.GetMsg(errercode.XRCode(code), GetLang(r))
		respData = failResp
	} else {
		respData = CreateSuccess(res)
	}
	resp, _ := json.MarshalIndent(respData, "", " ")
	r.Response.Write(resp)
	retTime := time.Now().Sub(startTime).Milliseconds()
	if err == nil {
		if retTime >= LongDoTime {
			g.Log().Infof(ctx, "reqId=%v,请求处理时间过长,处理时间=%vms,url=%v,ip=%v,authId=%v,响应数据=%v", r.GetHeader(ReqId, ""), retTime, r.URL.RequestURI(), r.GetClientIp(), userId, responseBodyForLog(r, resp))
		} else {
			g.Log().Infof(ctx, "reqId=%v,正常响应,处理时间=%vms,url=%v,ip=%v,authId=%v，响应数据=%v", r.GetHeader(ReqId, ""), retTime, r.URL.RequestURI(), r.GetClientIp(), userId, responseBodyForLog(r, resp))
		}
	}

}
